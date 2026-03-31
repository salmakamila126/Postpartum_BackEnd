package sleep

import (
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/errs"
	"sort"
)

func ValidateBulkInput(inputs []BulkInput, max int) error {
	if len(inputs) > max {
		return errs.ErrMaxSession
	}

	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].Start.Before(inputs[j].Start)
	})

	for i := 0; i < len(inputs); i++ {

		if inputs[i].End.Before(inputs[i].Start) {
			return errs.ErrInvalidTime
		}

		if i < len(inputs)-1 {
			if inputs[i].End.After(inputs[i+1].Start) {
				return errs.ErrOverlap
			}
		}
	}

	return nil
}

func ValidateCreateSession(
	sessions []entity.SleepSession,
	start, end int64,
	max int,
) error {

	for _, s := range sessions {

		if s.IsBackdate {
			return errs.ErrBackdateModify
		}

		if s.WakeTime == nil {
			continue
		}

		if start < s.WakeTime.Unix() && end > s.SleepTime.Unix() {
			return errs.ErrOverlap
		}
	}

	if len(sessions) >= max {
		return errs.ErrMaxSession
	}

	return nil
}
