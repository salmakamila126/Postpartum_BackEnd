package dto

type ScheduleSlot struct {
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Label     string `json:"label"`
}

type PsychologistListItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	Job          string `json:"job"`
	ExperienceYr int    `json:"experience_years"`
	PriceIDR     int    `json:"price_idr"`
	PhotoURL     string `json:"photo_url"`
}

type PsychologistDetail struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Title        string         `json:"title"`
	Job          string         `json:"job"`
	ExperienceYr int            `json:"experience_years"`
	PriceIDR     int            `json:"price_idr"`
	PhotoURL     string         `json:"photo_url"`
	Schedules    []ScheduleSlot `json:"schedules"`
}

type BookingWhatsAppRequest struct {
	SelectedSlot string `json:"selected_slot" binding:"required"`
}

type BookingWhatsAppResponse struct {
	WhatsAppURL    string `json:"whatsapp_url"`
	MessagePreview string `json:"message_preview"`
}

type UpdatePhotoURLRequest struct {
	PhotoURL string `json:"photo_url" binding:"required"`
}
