package usecase

import (
	"Postpartum_BackEnd/config"
	"Postpartum_BackEnd/internal/domain/sleep"
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/model"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/pkg/timeutil"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SleepUsecase struct {
	Repo   *repository.Repository
	DB     *gorm.DB
	Config config.SleepConfig
	Log    *zap.Logger
}

func NewSleepUsecase(repo *repository.Repository, db *gorm.DB, cfg config.SleepConfig, log *zap.Logger) *SleepUsecase {
	return &SleepUsecase{
		Repo:   repo,
		DB:     db,
		Config: cfg,
		Log:    log,
	}
}

func (u *SleepUsecase) logUser(userID uuid.UUID) *zap.Logger {
	return u.Log.With(zap.String("user_id", userID.String()))
}

func (u *SleepUsecase) withUserSleepTx(userID uuid.UUID, fn func(txRepo *repository.Repository) error) error {
	return u.DB.Transaction(func(tx *gorm.DB) error {
		var user entity.User
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		txRepo := repository.NewRepository(tx)
		return fn(txRepo)
	})
}

func (u *SleepUsecase) AddSleepSession(userID uuid.UUID, start, end time.Time) error {
	log := u.logUser(userID)
	log.Info("AddSleepSession called", zap.Time("start", start), zap.Time("end", end))

	now := timeutil.NowWIB()
	today := timeutil.StartOfDay(now)
	inputDate := timeutil.StartOfDay(start)

	if !inputDate.Equal(today) {
		log.Warn("manual input not today")
		return errs.ErrManualOnlyToday
	}

	if !end.After(start) {
		log.Warn("invalid time: end must be after start")
		return errs.ErrInvalidTime
	}

	date := timeutil.StartOfDay(start)
	if err := u.withUserSleepTx(userID, func(txRepo *repository.Repository) error {
		sessions, err := txRepo.SleepRepository.FindByDate(userID, date)
		if err != nil {
			return err
		}

		if err := sleep.ValidateCreateSession(sessions, start.Unix(), end.Unix(), u.Config.MaxSessionsPerDay); err != nil {
			return err
		}

		wakeTime := end
		return txRepo.SleepRepository.Create(&entity.SleepSession{
			UserID:     userID,
			SleepTime:  start,
			WakeTime:   &wakeTime,
			IsBackdate: false,
		})
	}); err != nil {
		log.Warn("manual sleep session transaction failed", zap.Error(err))
		return err
	}

	log.Info("sleep session created")
	return nil
}

func (u *SleepUsecase) StartSleep(userID uuid.UUID) error {
	log := u.logUser(userID)
	log.Info("StartSleep called")

	_, err := u.Repo.SleepRepository.FindActiveSession(userID)
	if err == nil {
		log.Warn("user already sleeping")
		return errs.ErrAlreadySleeping
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("failed to check active session", zap.Error(err))
		return err
	}

	now := timeutil.NowWIB()
	today := timeutil.StartOfDay(now)

	if err := u.withUserSleepTx(userID, func(txRepo *repository.Repository) error {
		_, err := txRepo.SleepRepository.FindActiveSession(userID)
		if err == nil {
			return errs.ErrAlreadySleeping
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		sessions, err := txRepo.SleepRepository.FindByDate(userID, today)
		if err != nil {
			return err
		}

		if err := sleep.ValidateCreateSession(sessions, now.Unix(), now.Unix(), u.Config.MaxSessionsPerDay); err != nil {
			return err
		}

		return txRepo.SleepRepository.Create(&entity.SleepSession{
			UserID:    userID,
			SleepTime: now,
		})
	}); err != nil {
		log.Warn("start sleep transaction failed", zap.Error(err))
		return err
	}

	log.Info("sleep session started", zap.Time("start_time", now))
	return nil
}

func (u *SleepUsecase) EndSleep(userID uuid.UUID) error {
	log := u.logUser(userID)
	log.Info("EndSleep called")

	now := timeutil.NowWIB()
	if err := u.withUserSleepTx(userID, func(txRepo *repository.Repository) error {
		session, err := txRepo.SleepRepository.FindActiveSession(userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrNoActiveSession
		}
		if err != nil {
			return err
		}

		session.WakeTime = &now
		return txRepo.SleepRepository.Update(session)
	}); err != nil {
		log.Warn("end sleep transaction failed", zap.Error(err))
		return err
	}

	log.Info("sleep session ended", zap.Time("end_time", now))
	return nil
}

func (u *SleepUsecase) GetDailySleep(userID uuid.UUID, date time.Time) (*dto.DailySleepResponse, error) {
	log := u.logUser(userID)
	log.Info("GetDailySleep called", zap.Time("date", date))

	sessions, err := u.Repo.SleepRepository.FindByDate(userID, date)
	if err != nil {
		log.Error("failed to fetch sessions", zap.Error(err))
		return nil, err
	}

	var totalMinutes int
	var count int
	var result []dto.SleepSessionItem
	isSleeping := false
	var currentStart string

	for _, s := range sessions {
		if s.WakeTime == nil {
			duration := int(timeutil.NowWIB().Sub(s.SleepTime).Minutes())
			isSleeping = true
			currentStart = timeutil.FormatHour(s.SleepTime)

			result = append(result, dto.SleepSessionItem{
				ID:              s.SleepID.String(),
				Start:           timeutil.FormatHour(s.SleepTime),
				End:             msgStillSleeping,
				DurationMinutes: duration,
				IsBackdate:      s.IsBackdate,
			})
			totalMinutes += duration
			count++
			continue
		}

		duration := int(s.WakeTime.Sub(s.SleepTime).Minutes())
		totalMinutes += duration
		count++

		result = append(result, dto.SleepSessionItem{
			ID:              s.SleepID.String(),
			Start:           timeutil.FormatHour(s.SleepTime),
			End:             timeutil.FormatHour(*s.WakeTime),
			DurationMinutes: duration,
			IsBackdate:      s.IsBackdate,
		})
	}

	if result == nil {
		result = []dto.SleepSessionItem{}
	}

	avg := 0
	if count > 0 {
		avg = totalMinutes / count
	}

	log.Info("daily sleep calculated", zap.Int("total_sessions", count))

	return &dto.DailySleepResponse{
		Date:              date.Format(timeutil.DateOnlyFormat),
		TotalSleepMinutes: totalMinutes,
		TotalSessions:     count,
		AvgSleepMinutes:   avg,
		Sessions:          result,
		IsSleeping:        isSleeping,
		CurrentStart:      currentStart,
	}, nil
}

func (u *SleepUsecase) GetHistory(userID uuid.UUID) (*dto.SleepHistoryResponse, error) {
	log := u.logUser(userID)
	log.Info("GetHistory called")

	sessions, err := u.Repo.SleepRepository.FindHistory(userID)
	if err != nil {
		log.Error("failed to fetch history", zap.Error(err))
		return nil, err
	}

	items := make([]dto.SleepHistoryItem, 0, len(sessions))
	for _, s := range sessions {
		item := dto.SleepHistoryItem{
			ID:         s.SleepID.String(),
			SleepTime:  timeutil.FormatHour(s.SleepTime),
			IsBackdate: s.IsBackdate,
			CreatedAt:  s.CreatedAt,
		}
		if s.WakeTime != nil {
			wakeStr := timeutil.FormatHour(*s.WakeTime)
			item.WakeTime = &wakeStr
			dur := int(s.WakeTime.Sub(s.SleepTime).Minutes())
			item.DurationMinutes = &dur
		}
		items = append(items, item)
	}

	return &dto.SleepHistoryResponse{Sessions: items}, nil
}

func (u *SleepUsecase) Predict(userID uuid.UUID) ([]model.SleepPrediction, error) {
	log := u.logUser(userID)
	log.Info("Predict called")

	now := timeutil.NowWIB()
	preds, err := u.predictFromAnchorDate(userID, now)
	if err != nil {
		log.Error("prediction failed", zap.Error(err))
		return nil, err
	}
	log.Info("prediction generated", zap.Int("count", len(preds)))

	return preds, nil
}

func (u *SleepUsecase) AddBulkSleepSession(userID uuid.UUID, date time.Time, inputs []dto.SleepManualRequest) error {
	log := u.logUser(userID)
	log.Info("AddBulkSleepSession called", zap.Int("count", len(inputs)))

	var bulkInputs []sleep.BulkInput

	for _, in := range inputs {
		start, err := timeutil.ParseRFC3339(in.Start)
		if err != nil {
			log.Warn("invalid start format", zap.Error(err))
			return errs.New(400, "invalid start time format: "+in.Start)
		}

		end, err := timeutil.ParseRFC3339(in.End)
		if err != nil {
			log.Warn("invalid end format", zap.Error(err))
			return errs.New(400, "invalid end time format: "+in.End)
		}

		if !end.After(start) {
			return errs.ErrInvalidTime
		}

		bulkInputs = append(bulkInputs, sleep.BulkInput{Start: start, End: end})
	}

	today := timeutil.StartOfDay(timeutil.NowWIB())
	targetDate := timeutil.StartOfDay(date)

	if !targetDate.Before(today) {
		log.Warn("bulk only past")
		return errs.ErrBulkOnlyPast
	}

	if err := sleep.ValidateBulkInput(bulkInputs, u.Config.MaxSessionsPerDay); err != nil {
		log.Warn("bulk validation failed", zap.Error(err))
		return err
	}

	var sessions []entity.SleepSession
	for _, in := range bulkInputs {
		endTime := in.End
		sessions = append(sessions, entity.SleepSession{
			UserID:     userID,
			SleepTime:  in.Start,
			WakeTime:   &endTime,
			IsBackdate: true,
		})
	}

	if err := u.withUserSleepTx(userID, func(txRepo *repository.Repository) error {
		existing, err := txRepo.SleepRepository.FindByDate(userID, targetDate)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return errs.ErrBulkAlreadyExist
		}

		return txRepo.SleepRepository.CreateBatch(sessions)
	}); err != nil {
		log.Warn("bulk transaction failed", zap.Error(err))
		return err
	}

	log.Info("bulk created", zap.Int("count", len(sessions)))
	return nil
}

func (u *SleepUsecase) GetStatus(userID uuid.UUID) (*dto.SleepStatusResponse, error) {
	log := u.logUser(userID)
	log.Info("GetStatus called")

	session, err := u.Repo.SleepRepository.FindActiveSession(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Info("user not sleeping")
		return &dto.SleepStatusResponse{Sleeping: false}, nil
	}
	if err != nil {
		log.Error("failed to get status", zap.Error(err))
		return nil, err
	}

	log.Info("user is sleeping")
	return &dto.SleepStatusResponse{
		Sleeping: true,
		Start:    timeutil.FormatHour(session.SleepTime),
	}, nil
}

func (u *SleepUsecase) GetTodayInsight(userID uuid.UUID) (*dto.InsightResponse, error) {
	log := u.logUser(userID)
	log.Info("GetTodayInsight called")

	now := timeutil.NowWIB()
	today := timeutil.StartOfDay(now)

	sessions, err := u.Repo.SleepRepository.FindByDate(userID, today)
	if err != nil {
		log.Error("failed to fetch today sessions", zap.Error(err))
		return nil, err
	}

	if len(sessions) == 0 {
		log.Info("no sleep data today")
		return &dto.InsightResponse{
			Status:  "empty",
			Message: msgEmpty,
		}, nil
	}

	last := sessions[len(sessions)-1]

	preds, err := u.predictFromAnchorDate(userID, today)
	if err != nil {
		if errors.Is(err, errs.ErrNotEnoughData) || errors.Is(err, errs.ErrNoData) {
			log.Warn("not enough data for insight", zap.Error(err))
			return &dto.InsightResponse{
				Status:  "not_enough_data",
				Message: msgNotEnoughData,
			}, nil
		}
		log.Error("prediction failed", zap.Error(err))
		return nil, err
	}
	if len(preds) == 0 {
		return nil, errs.ErrPredictionFailed
	}

	first := preds[0]

	if last.WakeTime == nil {
		log.Info("user currently sleeping", zap.Time("next_wake", first.NextWake))
		return &dto.InsightResponse{
			Status:  "sleeping",
			Message: msgWakeSoon + timeutil.FormatHour(first.NextWake) + msgWIB,
		}, nil
	}

	log.Info("user currently awake", zap.Time("next_sleep", first.NextSleep))
	return &dto.InsightResponse{
		Status:  "awake",
		Message: msgSleepSoon + timeutil.FormatHour(first.NextSleep) + msgWIB,
	}, nil
}

func (u *SleepUsecase) predictFromAnchorDate(userID uuid.UUID, anchorDate time.Time) ([]model.SleepPrediction, error) {
	sessions, err := u.Repo.SleepRepository.FindRecent(userID, u.Config.MaxRecentData)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, errs.ErrNoData
	}

	history := sleep.FilterPredictHistory(sessions, anchorDate)
	if len(history) == 0 {
		return nil, errs.ErrNoData
	}

	dayMap := sleep.GroupByDay(history)
	days := sleep.GetRecentDays(dayMap, u.Config.MinDaysForPredict)
	if len(days) < u.Config.MinDaysForPredict {
		return nil, errs.ErrNotEnoughData
	}

	avgSleep, avgWake, err := sleep.CalculateAverages(dayMap, days)
	if err != nil {
		return nil, err
	}

	anchorSessions := sleep.FilterValidSessions(sessions, timeutil.NowWIB())
	anchor, err := sleep.GetLatestFinishedSession(anchorSessions)
	if err != nil {
		return nil, errs.ErrNoData
	}

	anchorDay := timeutil.StartOfDay(anchor.SleepTime)
	targetDay := timeutil.StartOfDay(anchorDate)
	if anchorDay.After(targetDay) {
		return nil, errs.ErrNoData
	}
	if anchorDay.Equal(targetDay) {
		historyBeforeAnchor := sleep.FilterPredictHistory(sessions, anchorDay)
		dayMap = sleep.GroupByDay(historyBeforeAnchor)
		days = sleep.GetRecentDays(dayMap, u.Config.MinDaysForPredict)
		if len(days) < u.Config.MinDaysForPredict {
			return nil, errs.ErrNotEnoughData
		}
		avgSleep, avgWake, err = sleep.CalculateAverages(dayMap, days)
		if err != nil {
			return nil, err
		}
	}

	return sleep.GeneratePredictions(*anchor, avgSleep, avgWake, u.Config.PredictCycles), nil
}
