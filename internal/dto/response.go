package dto

type BaseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(data interface{}) BaseResponse {
	return BaseResponse{
		Success: true,
		Data:    data,
	}
}

func Error(message string) BaseResponse {
	return BaseResponse{
		Success: false,
		Error:   message,
	}
}
