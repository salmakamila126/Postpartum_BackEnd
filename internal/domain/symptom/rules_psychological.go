package symptom

func psychologicalCheckers() map[string]RuleChecker {
	return map[string]RuleChecker{
		"POSTPARTUM_DEPRESSION": func(i SymptomInput) (bool, []string) {
			if i.PPDTriggered {
				return true, []string{i.PPDDetail}
			}
			return false, nil
		},
	}
}
