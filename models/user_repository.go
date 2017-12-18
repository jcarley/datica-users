package models

import (
	"errors"

	"github.com/jcarley/datica-users/helper/encoding"
)

var (
	ErrMissingEmail          = errors.New("Missing email")
	ErrMissingUsername       = errors.New("Missing username")
	ErrMissingPassword       = errors.New("Missing password")
	ErrExistingUser          = errors.New("Existing user")
	ErrInvailUserCredentials = errors.New("Invalid user credentials")
)

type UserProvider interface {
	AddUser(user User) (User, error)
	FindByUsername(username string) (User, error)
	FindUser(userId string) (User, error)
}

type UserRepository struct {
	UserProvider
}

func NewUserRepository(userProvider UserProvider) *UserRepository {
	if userProvider == nil {
		panic("Nil UserProvider not allowed")
	}

	return &UserRepository{
		userProvider,
	}
}

func (this *UserRepository) SignIn(username, password string) (*User, error) {
	if username == "" {
		return &User{}, ErrMissingUsername
	}

	if password == "" {
		return &User{}, ErrMissingPassword
	}

	user, err := this.FindByUsername(username)
	if err != nil {
		return &User{}, err
	}

	enteredPassword := encoding.Key(password, user.Salt)
	if user.Password == enteredPassword {
		return &user, nil
	}
	return &User{}, ErrInvailUserCredentials
}
