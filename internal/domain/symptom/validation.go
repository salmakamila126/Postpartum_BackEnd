package symptom

import (
	"Postpartum_BackEnd/internal/dto"
	"errors"
	"strconv"
)

var ValidMoods = map[string]bool{
	"bahagia":      true,
	"tenang":       true,
	"bersyukur":    true,
	"bersemangat":  true,
	"percaya_diri": true,
	"optimis":      true,
	"sedih":        true,
	"cemas":        true,
	"mudah_marah":  true,
	"kewalahan":    true,
	"kesepian":     true,
	"putus_asa":    true,
}

func ValidateInput(req dto.CreateSymptomRequest) error {

	if len(req.Bleedings) == 0 {
		return errors.New("bleeding section is required (min 1 entry)")
	}
	if len(req.Bleedings) > 3 {
		return errors.New("maximum 3 bleeding entries allowed")
	}

	for i, b := range req.Bleedings {
		if b.PadUsage == "" {
			return errors.New("pad_usage is required for each bleeding entry")
		}
		if !ValidPads[b.PadUsage] {
			return errors.New("invalid pad_usage value at entry " + strconv.Itoa(i+1))
		}
		if b.ClotSize == "" {
			return errors.New("clot_size is required for each bleeding entry")
		}
		if !ValidClots[b.ClotSize] {
			return errors.New("invalid clot_size value at entry " + strconv.Itoa(i+1))
		}
		if b.BloodColor == "" {
			return errors.New("blood_color is required for each bleeding entry")
		}
		if !ValidColors[b.BloodColor] {
			return errors.New("invalid blood_color value at entry " + strconv.Itoa(i+1))
		}
		if b.Smell == "" {
			return errors.New("smell is required for each bleeding entry")
		}
		if !ValidSmells[b.Smell] {
			return errors.New("invalid smell value at entry " + strconv.Itoa(i+1))
		}
	}

	if !ValidTemperatures[req.Physical.Temperature] {
		return errors.New("invalid temperature value")
	}
	if err := validateStringList(req.Physical.Wound, ValidWounds, "wound"); err != nil {
		return err
	}
	if err := validateStringList(req.Physical.UrineProblems, ValidUrineProblems, "urine_problems"); err != nil {
		return err
	}
	if !ValidUrineColors[req.Physical.UrineColor] {
		return errors.New("invalid urine_color value")
	}
	if err := validateStringList(req.Physical.BreastProblems, ValidBreastProblems, "breast_problems"); err != nil {
		return err
	}
	if err := validateStringList(req.Physical.Swelling, ValidSwelling, "swelling"); err != nil {
		return err
	}
	if err := validateStringList(req.Physical.OtherSymptoms, ValidOtherSymptoms, "other_symptoms"); err != nil {
		return err
	}

	if len(req.Moods) == 0 {
		return errors.New("mood section is required (min 1 mood)")
	}
	if len(req.Moods) > 3 {
		return errors.New("maximum 3 moods allowed")
	}

	seen := make(map[string]bool)
	for _, m := range req.Moods {
		if !ValidMoods[m] {
			return errors.New("invalid mood value: " + m)
		}
		if seen[m] {
			return errors.New("duplicate mood value: " + m)
		}
		seen[m] = true
	}

	for _, scale := range []*int{
		req.Physical.Dizziness,
		req.Physical.Headache,
		req.Physical.Weakness,
		req.Physical.CalfPain,
		req.Physical.AbPain,
	} {
		if scale != nil && (*scale < 1 || *scale > 5) {
			return errors.New("scale values must be between 1 and 5")
		}
	}

	return nil
}

func validateStringList(values []string, allowed map[string]bool, field string) error {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if !allowed[value] {
			return errors.New("invalid " + field + " value: " + value)
		}
		if seen[value] {
			return errors.New("duplicate " + field + " value: " + value)
		}
		seen[value] = true
	}
	return nil
}
