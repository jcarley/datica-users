package models

type Token struct {
	TokenId string `json:"token_id,omitempty" db:"token_id"`
	Email   string `json:"email,omitempty" db:"email"`
	Token   string `json:"token,omitempty" db:"token"`
}
