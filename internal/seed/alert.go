package seed

import "Postpartum_BackEnd/internal/entity"

func AlertRuleSeedData() []entity.AlertRule {
	return []entity.AlertRule{
		{Code: "AB_PAIN", Level: "danger", Disease: "Infeksi atau komplikasi luka", Description: "Nyeri perut atau luka jahitan parah (level 5)", IsActive: true, IsSystem: true},
		{Code: "ANEMIA_DIZZINESS", Level: "danger", Disease: "Anemia", Description: "Pusing parah (level >= 4)", IsActive: true, IsSystem: true},
		{Code: "ANEMIA_WEAKNESS", Level: "danger", Disease: "Anemia", Description: "Lemas parah (level >= 4)", IsActive: true, IsSystem: true},
		{Code: "BLEEDING_INCREASE", Level: "warning", Disease: "Perdarahan pasca persalinan", Description: "Frekuensi ganti pembalut meningkat signifikan", IsActive: true, IsSystem: true},
		{Code: "BRIGHT_RED_LATE", Level: "danger", Disease: "Perdarahan abnormal pasca persalinan", Description: "Darah merah terang setelah minggu pertama", IsActive: true, IsSystem: true},
		{Code: "CALF_PAIN", Level: "danger", Disease: "Deep vein thrombosis", Description: "Nyeri betis parah (level 5)", IsActive: true, IsSystem: true},
		{Code: "CHEST_PAIN", Level: "danger", Disease: "Emboli paru", Description: "Nyeri dada", IsActive: true, IsSystem: true},
		{Code: "DARK_URINE", Level: "warning", Disease: "Dehidrasi atau infeksi", Description: "Warna urine gelap", IsActive: true, IsSystem: true},
		{Code: "FEVER", Level: "danger", Disease: "Infeksi", Description: "Suhu tubuh >= 38C", IsActive: true, IsSystem: true},
		{Code: "HEART_LUNG", Level: "danger", Disease: "Gangguan jantung atau paru", Description: "Nyeri dada dan sesak napas", IsActive: true, IsSystem: true},
		{Code: "HUGE_CLOT", Level: "danger", Disease: "Perdarahan pasca persalinan", Description: "Ukuran gumpalan: bola pingpong", IsActive: true, IsSystem: true},
		{Code: "LARGE_CLOT", Level: "warning", Disease: "Perdarahan pasca persalinan", Description: "Ukuran gumpalan: koin besar", IsActive: true, IsSystem: true},
		{Code: "MASTITIS", Level: "warning", Disease: "Mastitis (infeksi payudara)", Description: "Payudara bengkak, kemerahan, dan nyeri puting", IsActive: true, IsSystem: true},
		{Code: "NEURO_SEVERE", Level: "danger", Disease: "Gangguan saraf serius", Description: "Kejang, sakit kepala, dan pusing", IsActive: true, IsSystem: true},
		{Code: "POSTPARTUM_DEPRESSION", Level: "warning", Disease: "Kemungkinan Postpartum Depression", Description: "Emosi negatif mendominasi dalam pola mingguan", IsActive: true, IsSystem: true},
		{Code: "PPH", Level: "danger", Disease: "Perdarahan pasca persalinan", Description: "Pembalut penuh dalam kurang dari 2 jam", IsActive: true, IsSystem: true},
		{Code: "PREEKLAMPSIA", Level: "danger", Disease: "Preeklampsia pasca persalinan", Description: "Sakit kepala parah dengan tanda preeklampsia", IsActive: true, IsSystem: true},
		{Code: "SEIZURE", Level: "danger", Disease: "Eklampsia", Description: "Kejang-kejang", IsActive: true, IsSystem: true},
		{Code: "SHORTNESS_BREATH", Level: "danger", Disease: "Emboli paru", Description: "Sesak napas", IsActive: true, IsSystem: true},
		{Code: "SMELL", Level: "warning", Disease: "Infeksi", Description: "Cairan dari jalan lahir berbau menyengat", IsActive: true, IsSystem: true},
		{Code: "SWELLING_DANGER", Level: "danger", Disease: "Preeklampsia pasca persalinan", Description: "Pembengkakan di wajah atau tangan", IsActive: true, IsSystem: true},
		{Code: "URINE_CONTROL", Level: "warning", Disease: "Gangguan kandung kemih", Description: "Tidak bisa mengontrol BAK", IsActive: true, IsSystem: true},
		{Code: "URINE_RETENTION", Level: "warning", Disease: "Retensi urin", Description: "Tidak bisa BAK", IsActive: true, IsSystem: true},
		{Code: "UTI_FREQUENT", Level: "warning", Disease: "Infeksi saluran kemih", Description: "Ingin BAK terus menerus", IsActive: true, IsSystem: true},
		{Code: "UTI_PAIN", Level: "warning", Disease: "Infeksi saluran kemih", Description: "Nyeri saat BAK", IsActive: true, IsSystem: true},
		{Code: "WOUND_BLOOD", Level: "danger", Disease: "Komplikasi luka operasi", Description: "Perban luka ada bercak darah", IsActive: true, IsSystem: true},
		{Code: "WOUND_WET", Level: "danger", Disease: "Infeksi luka atau luka terbuka", Description: "Perban luka basah", IsActive: true, IsSystem: true},
	}
}
