package dto

type CreateSymptomRequest struct {
	Date      string            `json:"date" binding:"required"`
	Bleedings []BleedingDetail  `json:"bleedings"`
	Physical  PhysicalCondition `json:"physical"`
	Moods     []string          `json:"moods"`
}

type BleedingDetail struct {
	PadUsage   string `json:"pad_usage"`
	ClotSize   string `json:"clot_size"`
	BloodColor string `json:"blood_color"`
	Smell      string `json:"smell"`
}

type PhysicalCondition struct {
	Temperature string `json:"temperature"`

	Dizziness *int `json:"dizziness"`
	Headache  *int `json:"headache"`
	Weakness  *int `json:"weakness"`
	CalfPain  *int `json:"calf_pain"`
	AbPain    *int `json:"abdominal_pain"`

	Wound          []string `json:"wound"`
	UrineProblems  []string `json:"urine_problems"`
	UrineColor     string   `json:"urine_color"`
	BreastProblems []string `json:"breast_problems"`
	Swelling       []string `json:"swelling"`
	OtherSymptoms  []string `json:"other_symptoms"`
}

type AlertResponse struct {
	Level      string        `json:"level"`
	Confidence int           `json:"confidence"`
	Issues     []IssueResult `json:"issues"`
}

type IssueResult struct {
	Code        string   `json:"code"`
	Disease     string   `json:"disease"`
	Level       string   `json:"level"`
	Description string   `json:"description"`
	Symptoms    []string `json:"symptoms"`
}

type SymptomSaveResponse struct {
	Date  string        `json:"date"`
	Alert AlertResponse `json:"alert"`
}

type SymptomHistoryItem struct {
	ID         string            `json:"id"`
	Date       string            `json:"date"`
	IsBackdate bool              `json:"is_backdate"`
	Bleedings  []BleedingDetail  `json:"bleedings"`
	Physical   PhysicalCondition `json:"physical"`
	Moods      []string          `json:"moods"`
	Alert      AlertResponse     `json:"alert"`
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
}

type SymptomDetailResponse struct {
	SymptomHistoryItem
}
