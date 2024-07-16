package dto

type CreateUserRequest struct {
	FirstName string `json:"firstName"binding:"required"`
	LastName  string `json:"lastName"binding:"required"`
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
