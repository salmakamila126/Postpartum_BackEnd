package symptom

import "errors"

func DefaultCheckers() map[string]RuleChecker {
	checkers := make(map[string]RuleChecker)

	for code, checker := range bleedingCheckers() {
		checkers[code] = checker
	}
	for code, checker := range physicalCheckers() {
		checkers[code] = checker
	}
	for code, checker := range combinedCheckers() {
		checkers[code] = checker
	}
	for code, checker := range psychologicalCheckers() {
		checkers[code] = checker
	}

	return checkers
}

func BuildRules(defs []RuleDefinition) ([]Rule, error) {
	checkers := DefaultCheckers()
	rules := make([]Rule, 0, len(defs))

	for _, def := range defs {
		checker, ok := checkers[def.Code]
		if !ok {
			return nil, errors.New("missing checker for alert rule code: " + def.Code)
		}

		rules = append(rules, Rule{
			Code:        def.Code,
			Level:       def.Level,
			Disease:     def.Disease,
			Description: def.Description,
			Check:       checker,
		})
	}

	return rules, nil
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
