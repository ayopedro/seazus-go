package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,min=6"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"-" validate:"required,min=6"`
}

type AuthUser struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type AuthResponse struct {
	User  AuthUser `json:"user"`
	Token string   `json:"access_token"`
}
