package symptom

type SymptomInput struct {
	Bleedings      []BleedingInput
	Physical       PhysicalInput
	Moods          []string
	MoodHistory    [][]string
	DaysAfterBirth int

	PPDTriggered bool
	PPDDetail    string
}

type BleedingInput struct {
	PadUsage   string
	ClotSize   string
	BloodColor string
	Smell      string
}

type PhysicalInput struct {
	Temperature string

	Dizziness int
	Headache  int
	Weakness  int
	CalfPain  int
	AbPain    int

	Wound          []string
	UrineProblems  []string
	UrineColor     string
	BreastProblems []string
	Swelling       []string
	OtherSymptoms  []string
}
