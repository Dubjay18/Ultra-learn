package errors

const (
	ValidationError     = "Bad request"
	InternalServerError = "Something went wrong"
	UnAuthorized        = "Unauthorized"
)

type ApiError struct {
	StatusCode int    `json:"statusCode"`
	Error      any    `json:"error"`
	Message    string `json:"message"`
}
