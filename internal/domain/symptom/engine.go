package symptom

type AlertResult struct {
	Level      AlertLevel
	Issues     []Issue
	Confidence int
}

type Issue struct {
	Code        string
	Disease     string
	Level       AlertLevel
	Description string
	Symptoms    []string
}

func Evaluate(input SymptomInput, rules []Rule) AlertResult {
	var issues []Issue

	for _, rule := range rules {
		ok, symptoms := rule.Check(input)
		if !ok {
			continue
		}
		issues = append(issues, Issue{
			Code:        rule.Code,
			Disease:     rule.Disease,
			Level:       rule.Level,
			Description: rule.Description,
			Symptoms:    symptoms,
		})
	}

	if len(issues) == 0 {
		return AlertResult{
			Level: Safe,
			Issues: []Issue{
				{
					Code:        "NORMAL",
					Disease:     "Kondisi normal",
					Level:       Safe,
					Description: "Tidak ada tanda bahaya yang terdeteksi",
				},
			},
			Confidence: 100,
		}
	}

	priority := map[AlertLevel]int{Safe: 0, Warning: 1, Danger: 2}
	highestLevel := Safe
	for _, iss := range issues {
		if priority[iss.Level] > priority[highestLevel] {
			highestLevel = iss.Level
		}
	}

	var topIssue *Issue
	for _, iss := range issues {
		if iss.Level != highestLevel {
			continue
		}

		if topIssue == nil || isBetterIssue(iss, *topIssue) {
			current := iss
			topIssue = &current
		}
	}

	if topIssue == nil {
		return AlertResult{
			Level:      highestLevel,
			Issues:     []Issue{},
			Confidence: 25,
		}
	}

	confidence := len(topIssue.Symptoms) * ConfidencePerSymptom
	if confidence > 100 {
		confidence = 100
	}
	if confidence < 25 {
		confidence = 25
	}

	return AlertResult{
		Level:      highestLevel,
		Issues:     []Issue{*topIssue},
		Confidence: confidence,
	}
}

func isBetterIssue(candidate Issue, current Issue) bool {
	if len(candidate.Symptoms) != len(current.Symptoms) {
		return len(candidate.Symptoms) > len(current.Symptoms)
	}

	return candidate.Code < current.Code
}
