package dto

type ApiSuccessResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Data       any    `json:"data"`
}
