package dto

import "time"

// RegisterRequest holds user registration payload.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=driver admin"`
}

// LoginRequest holds user login credentials.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user profile data without password.
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginUserResponse is the nested user info in login response.
type LoginUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// LoginResponse represents successful login payload.
type LoginResponse struct {
	Token string            `json:"token"`
	User  LoginUserResponse `json:"user"`
}
