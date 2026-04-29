package models

import "time"

type CreateUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type LoginUser struct {
	Email    string
	Password string
}

type User struct {
	Id         string
	FirstName  string
	LastName   string
	Email      string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsVerified bool
}

type AuthUser struct {
	Id         string
	FirstName  string
	Email      string
	IsVerified bool
}

type AuthResponse struct {
	User  AuthUser
	Token string
}
