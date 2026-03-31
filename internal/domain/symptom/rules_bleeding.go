package symptom

func padToMinutes(p string) int {
	switch p {
	case "24h":
		return 1440
	case "6h":
		return 360
	case "2h":
		return 120
	case "<2h":
		return 90
	default:
		return 0
	}
}

func isDangerPadUsage(p string) bool {
	return p == "<2h"
}

func bleedingCheckers() map[string]RuleChecker {
	return map[string]RuleChecker{
		"PPH": func(i SymptomInput) (bool, []string) {
			for _, b := range i.Bleedings {
				if isDangerPadUsage(b.PadUsage) {
					return true, []string{"Pembalut penuh dalam kurang dari 2 jam"}
				}
			}
			return false, nil
		},
		"HUGE_CLOT": func(i SymptomInput) (bool, []string) {
			for _, b := range i.Bleedings {
				if b.ClotSize == "pingpong" {
					return true, []string{"Gumpalan darah ukuran bola pingpong"}
				}
			}
			return false, nil
		},
		"LARGE_CLOT": func(i SymptomInput) (bool, []string) {
			for _, b := range i.Bleedings {
				if b.ClotSize == "large_coin" {
					return true, []string{"Gumpalan darah ukuran koin besar"}
				}
			}
			return false, nil
		},
		"BRIGHT_RED_LATE": func(i SymptomInput) (bool, []string) {
			if i.DaysAfterBirth <= 7 {
				return false, nil
			}
			for _, b := range i.Bleedings {
				if b.BloodColor == "bright_red" {
					return true, []string{"Darah merah terang setelah minggu pertama pasca persalinan"}
				}
			}
			return false, nil
		},
		"SMELL": func(i SymptomInput) (bool, []string) {
			for _, b := range i.Bleedings {
				if b.Smell == "strong" {
					return true, []string{"Cairan dari jalan lahir berbau menyengat"}
				}
			}
			return false, nil
		},
		"BLEEDING_INCREASE": func(i SymptomInput) (bool, []string) {
			if len(i.Bleedings) < 2 {
				return false, nil
			}
			for idx := 1; idx < len(i.Bleedings); idx++ {
				prev := padToMinutes(i.Bleedings[idx-1].PadUsage)
				curr := padToMinutes(i.Bleedings[idx].PadUsage)
				if curr > 0 && prev > 0 && curr <= prev/2 {
					return true, []string{"Frekuensi ganti pembalut meningkat 2x dari entri sebelumnya"}
				}
			}
			return false, nil
		},
	}
}
