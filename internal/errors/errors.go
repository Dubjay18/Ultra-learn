package errors

const (
	ValidationError     = "Bad request"
	InternalServerError = "Something went wrong"
	UnAuthorized        = "Unauthorized"
	UserNotFound        = "User not found"
)

type ApiError struct {
	StatusCode int    `json:"status_code"`
	Error      any    `json:"error"`
	Message    string `json:"message"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
