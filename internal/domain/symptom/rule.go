package symptom

type Rule struct {
	Code        string
	Level       AlertLevel
	Disease     string
	Description string
	Check       func(input SymptomInput) (bool, []string)
}

type RuleChecker func(input SymptomInput) (bool, []string)

type RuleDefinition struct {
	Code        string
	Level       AlertLevel
	Disease     string
	Description string
}
