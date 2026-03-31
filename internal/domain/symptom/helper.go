package symptom

import (
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/pkg/timeutil"
	"sort"
	"strconv"
	"time"
)

var NegativeMoods = map[string]bool{
	"sedih":       true,
	"cemas":       true,
	"mudah_marah": true,
	"kewalahan":   true,
	"kesepian":    true,
	"putus_asa":   true,
}

type WeeklyMoodResult struct {
	WeekNumber    int
	NegDaysCount  int
	TotalDaysData int
	Triggered     bool
}

func isDayNegativeDominant(moods []string) bool {
	negCount := 0
	for _, m := range moods {
		if NegativeMoods[m] {
			negCount++
		}
	}
	return negCount >= 1
}

func AnalyzePPDWeekly(logs []entity.SymptomLog) (triggered bool, detail string) {
	if len(logs) == 0 {
		return false, ""
	}

	weekMap := make(map[time.Time][]entity.SymptomLog)
	for _, l := range logs {
		monday := startOfWeek(l.Date)
		weekMap[monday] = append(weekMap[monday], l)
	}

	var weeks []time.Time
	for w := range weekMap {
		weeks = append(weeks, w)
	}
	sortTimes(weeks)

	type weekStat struct {
		weekStart time.Time
		negDays   int
		totalDays int
	}

	var stats []weekStat
	for _, w := range weeks {
		logsInWeek := weekMap[w]
		negDays := 0
		for _, l := range logsInWeek {
			var moods []string
			for _, m := range l.Moods {
				moods = append(moods, m.Type)
			}
			if isDayNegativeDominant(moods) {
				negDays++
			}
		}
		stats = append(stats, weekStat{
			weekStart: w,
			negDays:   negDays,
			totalDays: len(logsInWeek),
		})
	}

	for _, s := range stats {
		if s.negDays >= 4 {
			return true, "Emosi negatif mendominasi " + strconv.Itoa(s.negDays) + " hari dalam minggu " + s.weekStart.Format("02 Jan 2006")
		}
	}

	if len(stats) >= 3 {
		for i := 0; i+2 < len(stats); i++ {
			w1 := stats[i]
			w2 := stats[i+1]
			w3 := stats[i+2]

			if !w2.weekStart.Equal(w1.weekStart.AddDate(0, 0, 7)) {
				continue
			}
			if !w3.weekStart.Equal(w2.weekStart.AddDate(0, 0, 7)) {
				continue
			}

			if w1.negDays >= 3 && w2.negDays >= 3 && w3.negDays >= 1 {
				return true, "Pola emosi negatif berlanjut: " +
					strconv.Itoa(w1.negDays) + " hari (minggu 1), " +
					strconv.Itoa(w2.negDays) + " hari (minggu 2), dan masih berlanjut di minggu ke-3"
			}
		}
	}

	if len(stats) >= 2 {
		for i := 0; i+1 < len(stats); i++ {
			w1 := stats[i]
			w2 := stats[i+1]
			if !w2.weekStart.Equal(w1.weekStart.AddDate(0, 0, 7)) {
				continue
			}
			if w1.negDays >= 3 && w2.negDays >= 3 {
				return true, "Emosi negatif mendominasi " +
					strconv.Itoa(w1.negDays) + " hari dan " +
					strconv.Itoa(w2.negDays) + " hari dalam 2 minggu berturut-turut"
			}
		}
	}

	return false, ""
}

func startOfWeek(t time.Time) time.Time {
	t = timeutil.StartOfDay(t)

	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	return t.AddDate(0, 0, -(weekday - 1))
}

func sortTimes(ts []time.Time) {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Before(ts[j])
	})
}
