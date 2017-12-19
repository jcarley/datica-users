package models

import (
	"errors"
)

var (
	ErrMissingToken  = errors.New("Missing token")
	ErrTokenNotFound = errors.New("Token not found")
)

type TokenProvider interface {
	AddToken(token Token) (Token, error)
	FindToken(token string) (Token, error)
}

type TokenRepository struct {
	TokenProvider
}

func NewTokenRepository(tokenProvider TokenProvider) *TokenRepository {
	if tokenProvider == nil {
		panic("Nil TokenProvider not allowed")
	}

	return &TokenRepository{
		tokenProvider,
	}
}

func (this *TokenRepository) Exists(token string) (bool, error) {

	if token == "" {
		return false, ErrMissingToken
	}

	_, err := this.FindToken(token)
	if err != nil {
		return false, err
	}

	return true, nil
}
