package sleep

import (
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/model"
	"Postpartum_BackEnd/pkg/timeutil"
	"errors"
	"sort"
	"time"
)

func GroupByDay(sessions []entity.SleepSession) map[string][]entity.SleepSession {
	result := make(map[string][]entity.SleepSession)

	for _, s := range sessions {
		day := s.SleepTime.Format(timeutil.DateOnlyFormat)
		result[day] = append(result[day], s)
	}

	return result
}

func GetRecentDays(dayMap map[string][]entity.SleepSession, limit int) []string {
	var days []string

	for d := range dayMap {
		days = append(days, d)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(days)))

	if len(days) < limit {
		return days
	}

	return days[:limit]
}

func CalculateAverages(
	dayMap map[string][]entity.SleepSession,
	days []string,
) (time.Duration, time.Duration, error) {

	var totalSleep time.Duration
	var totalWake time.Duration
	var countSleep int
	var countWake int

	for _, d := range days {
		sessions := dayMap[d]

		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].SleepTime.Before(sessions[j].SleepTime)
		})

		for i := 0; i < len(sessions); i++ {
			s := sessions[i]

			if s.WakeTime == nil {
				continue
			}

			totalSleep += s.WakeTime.Sub(s.SleepTime)
			countSleep++

			if i < len(sessions)-1 {
				next := sessions[i+1]

				if next.SleepTime.After(*s.WakeTime) {
					totalWake += next.SleepTime.Sub(*s.WakeTime)
					countWake++
				}
			}
		}
	}

	if countSleep == 0 || countWake == 0 {
		return 0, 0, errors.New("not enough data to calculate averages")
	}

	return totalSleep / time.Duration(countSleep),
		totalWake / time.Duration(countWake),
		nil
}

func GeneratePredictions(
	last entity.SleepSession,
	avgSleep time.Duration,
	avgWake time.Duration,
	cycles int,
) []model.SleepPrediction {

	nextSleep := last.WakeTime.Add(avgWake)
	nextWake := nextSleep.Add(avgSleep)

	var results []model.SleepPrediction

	currentSleep := nextSleep
	currentWake := nextWake

	for i := 0; i < cycles; i++ {
		results = append(results, model.SleepPrediction{
			NextSleep: currentSleep,
			NextWake:  currentWake,
		})

		currentSleep = currentWake.Add(avgWake)
		currentWake = currentSleep.Add(avgSleep)
	}

	return results
}

func FilterValidSessions(sessions []entity.SleepSession, now time.Time) []entity.SleepSession {
	var result []entity.SleepSession

	for _, s := range sessions {
		if s.SleepTime.Before(now) && !s.IsBackdate {
			result = append(result, s)
		}
	}

	return result
}

func FilterPredictHistory(sessions []entity.SleepSession, before time.Time) []entity.SleepSession {
	var result []entity.SleepSession

	cutoff := timeutil.StartOfDay(before)
	for _, s := range sessions {
		if s.IsBackdate || s.WakeTime == nil {
			continue
		}
		if timeutil.StartOfDay(s.SleepTime).Before(cutoff) {
			result = append(result, s)
		}
	}

	return result
}

func GetLatestFinishedSession(sessions []entity.SleepSession) (*entity.SleepSession, error) {
	var latest *entity.SleepSession

	for i := range sessions {
		s := sessions[i]
		if s.IsBackdate || s.WakeTime == nil {
			continue
		}
		if latest == nil || s.SleepTime.After(latest.SleepTime) {
			copy := s
			latest = &copy
		}
	}

	if latest == nil {
		return nil, errors.New("no finished session")
	}

	return latest, nil
}
