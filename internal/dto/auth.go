package dto

type CreateUserRequest struct {
	FirstName string `json:"first_name"binding:"required"`
	LastName  string `json:"last_name"binding:"required"`
	Email     string `json:"email"binding:"required""`
	Password  string `json:"password"binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"binding:"required"`
	Password string `json:"password"binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
