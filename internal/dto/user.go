package dto

type UserResponse struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}
