package models

import "time"

type User struct {
	UserId    string    `json:"user_id,omitempty" db:"user_id"`
	Name      string    `json:"name,omitempty" db:"name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Salt      string    `json:"salt,omitempty" db:"salt"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
