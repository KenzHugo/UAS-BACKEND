package model

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	RoleID       *string   `json:"role_id" db:"role_id"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Role struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Permission struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Resource    string  `json:"resource" db:"resource"`
	Action      string  `json:"action" db:"action"`
	Description *string `json:"description" db:"description"`
}

type UserWithRole struct {
	User
	RoleName    string   `json:"role" db:"role_name"`
	Permissions []string `json:"permissions"`
}

// Request/Response DTOs
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	User         UserProfile `json:"user"`
}

type UserProfile struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"fullName"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}