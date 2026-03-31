package symptom

import "strconv"

func combinedCheckers() map[string]RuleChecker {
	return map[string]RuleChecker{
		"NEURO_SEVERE": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.OtherSymptoms, "kejang") &&
				hasHeadache(i.Physical) &&
				i.Physical.Dizziness >= 1 {
				return true, []string{"Kejang-kejang", "Sakit kepala", "Pusing"}
			}
			return false, nil
		},
		"HEART_LUNG": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.OtherSymptoms, "nyeri_dada") &&
				contains(i.Physical.OtherSymptoms, "sesak_napas") {
				return true, []string{"Nyeri dada", "Sesak napas"}
			}
			return false, nil
		},
		"PREEKLAMPSIA": func(i SymptomInput) (bool, []string) {
			if headacheLevel(i.Physical) >= 4 {
				var triggered []string
				if contains(i.Physical.OtherSymptoms, "penglihatan_kabur") {
					triggered = append(triggered, "Penglihatan kabur")
				}
				if contains(i.Physical.OtherSymptoms, "muntah") {
					triggered = append(triggered, "Muntah")
				}
				if contains(i.Physical.OtherSymptoms, "nyeri_ulu_hati") {
					triggered = append(triggered, "Nyeri ulu hati")
				}
				if contains(i.Physical.Swelling, "kaki") {
					triggered = append(triggered, "Bengkak pergelangan kaki")
				}
				if contains(i.Physical.Swelling, "wajah") {
					triggered = append(triggered, "Bengkak wajah")
				}
				if contains(i.Physical.Swelling, "tangan") {
					triggered = append(triggered, "Bengkak tangan")
				}
				if len(triggered) > 0 {
					all := append([]string{"Sakit kepala parah (level " + strconv.Itoa(headacheLevel(i.Physical)) + ")"}, triggered...)
					return true, all
				}
			}
			return false, nil
		},
	}
}

func hasHeadache(physical PhysicalInput) bool {
	return physical.Headache >= 1 || contains(physical.OtherSymptoms, "sakit_kepala")
}

func headacheLevel(physical PhysicalInput) int {
	if physical.Headache > 0 {
		return physical.Headache
	}
	if contains(physical.OtherSymptoms, "sakit_kepala") {
		return 1
	}
	return 0
}
