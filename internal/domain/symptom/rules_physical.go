package symptom

import "strconv"

func physicalCheckers() map[string]RuleChecker {
	return map[string]RuleChecker{
		"FEVER": func(i SymptomInput) (bool, []string) {
			if i.Physical.Temperature == ">=38" {
				return true, []string{"Demam tinggi (>= 38C)"}
			}
			return false, nil
		},
		"ANEMIA_DIZZINESS": func(i SymptomInput) (bool, []string) {
			if i.Physical.Dizziness >= 4 {
				return true, []string{"Pusing level " + strconv.Itoa(i.Physical.Dizziness)}
			}
			return false, nil
		},
		"ANEMIA_WEAKNESS": func(i SymptomInput) (bool, []string) {
			if i.Physical.Weakness >= 4 {
				return true, []string{"Lemas level " + strconv.Itoa(i.Physical.Weakness)}
			}
			return false, nil
		},
		"SEIZURE": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.OtherSymptoms, "kejang") {
				return true, []string{"Kejang-kejang"}
			}
			return false, nil
		},
		"SWELLING_DANGER": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.Swelling, "wajah") || contains(i.Physical.Swelling, "tangan") {
				return true, []string{"Pembengkakan tidak normal (wajah/tangan)"}
			}
			return false, nil
		},
		"CHEST_PAIN": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.OtherSymptoms, "nyeri_dada") {
				return true, []string{"Nyeri dada"}
			}
			return false, nil
		},
		"SHORTNESS_BREATH": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.OtherSymptoms, "sesak_napas") {
				return true, []string{"Sesak napas"}
			}
			return false, nil
		},
		"CALF_PAIN": func(i SymptomInput) (bool, []string) {
			if i.Physical.CalfPain >= 5 {
				return true, []string{"Nyeri betis level " + strconv.Itoa(i.Physical.CalfPain)}
			}
			return false, nil
		},
		"AB_PAIN": func(i SymptomInput) (bool, []string) {
			if i.Physical.AbPain >= 5 {
				return true, []string{"Nyeri perut/luka jahitan level " + strconv.Itoa(i.Physical.AbPain)}
			}
			return false, nil
		},
		"WOUND_BLOOD": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.Wound, "bercak_darah") {
				return true, []string{"Bercak darah pada perban luka"}
			}
			return false, nil
		},
		"WOUND_WET": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.Wound, "basah") {
				return true, []string{"Perban luka basah"}
			}
			return false, nil
		},
		"UTI_PAIN": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.UrineProblems, "nyeri_bak") {
				return true, []string{"Nyeri saat BAK"}
			}
			return false, nil
		},
		"UTI_FREQUENT": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.UrineProblems, "sering_bak") {
				return true, []string{"Ingin BAK terus menerus"}
			}
			return false, nil
		},
		"DARK_URINE": func(i SymptomInput) (bool, []string) {
			if i.Physical.UrineColor == "dark" {
				return true, []string{"Warna urine gelap"}
			}
			return false, nil
		},
		"URINE_RETENTION": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.UrineProblems, "tidak_bisa_bak") {
				return true, []string{"Tidak bisa BAK"}
			}
			return false, nil
		},
		"URINE_CONTROL": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.UrineProblems, "tidak_kontrol") {
				return true, []string{"Tidak bisa mengontrol BAK"}
			}
			return false, nil
		},
		"MASTITIS": func(i SymptomInput) (bool, []string) {
			if contains(i.Physical.BreastProblems, "bengkak") &&
				contains(i.Physical.BreastProblems, "kemerahan") &&
				contains(i.Physical.BreastProblems, "nyeri_puting") {
				return true, []string{"Payudara bengkak", "Payudara kemerahan", "Nyeri puting"}
			}
			return false, nil
		},
	}
}
