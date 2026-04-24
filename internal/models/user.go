package models

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id         string   `json:"id"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Password   string   `json:"-"`
	CreatedAt  JSONTime `json:"created_at,omitempty"`
	UpdatedAt  JSONTime `json:"updated_at,omitempty"`
	IsVerified bool     `json:"is_verified"`
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
