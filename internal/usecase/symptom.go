package usecase

import (
	"Postpartum_BackEnd/internal/domain/symptom"
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/internal/seed"
	"Postpartum_BackEnd/pkg/cache"
	"Postpartum_BackEnd/pkg/timeutil"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SymptomUsecase struct {
	Repo  *repository.Repository
	Cache *cache.Cache
}

func NewSymptomUsecase(repo *repository.Repository, appCache *cache.Cache) *SymptomUsecase {
	return &SymptomUsecase{Repo: repo, Cache: appCache}
}

func (u *SymptomUsecase) SeedAlertRulesIfEmpty() error {
	count, err := u.Repo.AlertRepository.CountAll()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if err := u.Repo.AlertRepository.SeedAll(seed.AlertRuleSeedData()); err != nil {
		return err
	}
	if err := u.Cache.Delete(context.Background(), "alert_rules:active"); err != nil {
		log.Printf("cache delete failed for alert rules: %v", err)
	}
	return nil
}

func (u *SymptomUsecase) CreateOrUpdate(
	userID uuid.UUID,
	req dto.CreateSymptomRequest,
) (*symptom.AlertResult, error) {
	if err := symptom.ValidateInput(req); err != nil {
		return nil, errs.New(http.StatusBadRequest, err.Error())
	}

	date, err := timeutil.ParseDate(req.Date)
	if err != nil {
		return nil, errs.New(http.StatusBadRequest, "invalid date format - use YYYY-MM-DD")
	}

	today := timeutil.StartOfDay(timeutil.NowWIB())
	target := timeutil.StartOfDay(date)

	if target.After(today) {
		return nil, errs.ErrFutureDate
	}

	isToday := target.Equal(today)

	existing, findErr := u.Repo.SymptomRepository.FindByDate(userID, date)
	hasExisting := findErr == nil && existing != nil && existing.ID != uuid.Nil
	if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, findErr
	}

	logEntity, err := mapToEntity(userID, req)
	if err != nil {
		return nil, err
	}
	logEntity.Date = date

	moodLogs, moodErr := u.Repo.SymptomRepository.FindMoodHistoryWeekly(userID, 4)
	if moodErr != nil {
		moodLogs = []entity.SymptomLog{}
	}
	moodLogs = mergeCurrentMoodLog(moodLogs, logEntity)
	ppdTriggered, ppdDetail := symptom.AnalyzePPDWeekly(moodLogs)

	baby, err := u.Repo.BabyRepository.FindByUserID(userID)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, "baby record not found for this user")
	}
	birthDate, err := parseStoredBirthDate(baby.BirthDate)
	if err != nil {
		return nil, errs.New(http.StatusBadRequest, "stored birth date format is invalid; update profile birth_date using YYYY-MM-DD")
	}
	daysAfterBirth := int(timeutil.NowWIB().Sub(birthDate).Hours() / 24)

	input := toDomainInput(req)
	input.DaysAfterBirth = daysAfterBirth
	input.PPDTriggered = ppdTriggered
	input.PPDDetail = ppdDetail

	ruleEntities, err := u.findActiveAlertRules()
	if err != nil {
		return nil, err
	}

	defs := make([]symptom.RuleDefinition, 0, len(ruleEntities))
	for _, rule := range ruleEntities {
		defs = append(defs, symptom.RuleDefinition{
			Code:        rule.Code,
			Level:       symptom.AlertLevel(rule.Level),
			Disease:     rule.Disease,
			Description: rule.Description,
		})
	}

	rules, err := symptom.BuildRules(defs)
	if err != nil {
		return nil, err
	}

	alert := symptom.Evaluate(input, rules)

	if err := applyAlertSnapshot(logEntity, &alert); err != nil {
		return nil, err
	}

	if err := u.saveDailyLog(logEntity, existing, hasExisting, isToday); err != nil {
		return nil, err
	}

	return &alert, nil
}

func (u *SymptomUsecase) saveDailyLog(
	logEntity *entity.SymptomLog,
	existing *entity.SymptomLog,
	hasExisting bool,
	isToday bool,
) error {
	if isToday {
		if hasExisting {
			logEntity.ID = existing.ID
			logEntity.IsBackdate = false
			return u.Repo.SymptomRepository.Update(logEntity)
		}

		logEntity.IsBackdate = false
		return u.Repo.SymptomRepository.Create(logEntity)
	}

	if hasExisting {
		return errs.ErrBackdateExists
	}

	logEntity.IsBackdate = true
	return u.Repo.SymptomRepository.Create(logEntity)
}

func mergeCurrentMoodLog(logs []entity.SymptomLog, current *entity.SymptomLog) []entity.SymptomLog {
	for i := range logs {
		if timeutil.StartOfDay(logs[i].Date).Equal(timeutil.StartOfDay(current.Date)) {
			logs[i].Moods = current.Moods
			return logs
		}
	}

	return append(logs, *current)
}

func parseStoredBirthDate(value string) (time.Time, error) {
	layouts := []string{
		timeutil.DateOnlyFormat,
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, timeutil.NowWIB().Location()); err == nil {
			return t, nil
		}
	}

	if len(value) >= len(timeutil.DateOnlyFormat) {
		if t, err := timeutil.ParseDate(value[:len(timeutil.DateOnlyFormat)]); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errs.New(http.StatusBadRequest, "invalid stored birth date")
}

func (u *SymptomUsecase) GetHistory(userID uuid.UUID) ([]dto.SymptomHistoryItem, error) {
	logs, err := u.Repo.SymptomRepository.FindHistory(userID)
	if err != nil {
		return nil, err
	}
	return toHistoryResponse(logs)
}

func (u *SymptomUsecase) GetDetail(userID uuid.UUID, dateStr string) (*dto.SymptomDetailResponse, error) {
	date, err := timeutil.ParseDate(dateStr)
	if err != nil {
		return nil, errs.New(http.StatusBadRequest, "invalid date format - use YYYY-MM-DD")
	}

	data, err := u.Repo.SymptomRepository.FindByDate(userID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSymptomNotFound
		}
		return nil, err
	}

	return toDetailResponse(data)
}

func (u *SymptomUsecase) findActiveAlertRules() ([]entity.AlertRule, error) {
	ctx := context.Background()
	const cacheKey = "alert_rules:active"

	var cached []entity.AlertRule
	if ok, err := u.Cache.GetJSON(ctx, cacheKey, &cached); err == nil && ok {
		return cached, nil
	}

	rules, err := u.Repo.AlertRepository.FindActive()
	if err != nil {
		return nil, err
	}

	if err := u.Cache.SetJSON(ctx, cacheKey, rules, 30*time.Minute); err != nil {
		log.Printf("cache set failed for %s: %v", cacheKey, err)
	}
	return rules, nil
}
