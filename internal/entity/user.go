package entity

import "time"

type User struct {
	ID           int       `json:"id,omitempty"`
	Email        string    `json:"email,omitempty"`
	Login        string    `json:"login,omitempty"`
	Password     string    `json:"password,omitempty"`
	RegisteredAt time.Time `json:"registered_at,omitempty"`
}
