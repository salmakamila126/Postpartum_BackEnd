package symptom

type AlertLevel string

const (
	Safe    AlertLevel = "safe"
	Warning AlertLevel = "warning"
	Danger  AlertLevel = "danger"
)

const (
	Pad24h  = "24h"
	Pad6h   = "6h"
	Pad2h   = "2h"
	PadFast = "<2h"
)

var ValidPads = map[string]bool{
	Pad24h:  true,
	Pad6h:   true,
	Pad2h:   true,
	PadFast: true,
}

var ValidClots = map[string]bool{
	"none":       true,
	"small_coin": true,
	"large_coin": true,
	"pingpong":   true,
}

var ValidColors = map[string]bool{
	"dark_red":   true,
	"normal_red": true,
	"bright_red": true,
}

var ValidSmells = map[string]bool{
	"none":   true,
	"mild":   true,
	"strong": true,
}

var ValidTemperatures = map[string]bool{
	"":     true,
	"<=36": true,
	"36.5": true,
	"37":   true,
	"37.5": true,
	">=38": true,
}

var ValidWounds = map[string]bool{
	"tidak_ada_perban": true,
	"kering":           true,
	"bercak_darah":     true,
	"basah":            true,
}

var ValidUrineProblems = map[string]bool{
	"tidak_bisa_bak": true,
	"tidak_kontrol":  true,
	"sering_bak":     true,
	"nyeri_bak":      true,
}

var ValidUrineColors = map[string]bool{
	"":     true,
	"dark": true,
}

var ValidBreastProblems = map[string]bool{
	"bengkak":        true,
	"kemerahan":      true,
	"nyeri_puting":   true,
	"nyeri_payudara": true,
}

var ValidSwelling = map[string]bool{
	"kaki":   true,
	"tangan": true,
	"wajah":  true,
}

var ValidOtherSymptoms = map[string]bool{
	"kejang":            true,
	"sakit_kepala":      true,
	"penglihatan_kabur": true,
	"muntah":            true,
	"nyeri_ulu_hati":    true,
	"nyeri_dada":        true,
	"sesak_napas":       true,
}

const ConfidencePerSymptom = 25
