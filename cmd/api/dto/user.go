package dto

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
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
