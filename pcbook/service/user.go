package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string
	HashedPassword []byte
	Role           string
}

func NewUser(username string, password string, role string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("invalid username or password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &User{
		Username:       username,
		HashedPassword: hashedPassword,
		Role:           role,
	}, nil
}

func (user *User) IsPasswordValid(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	return err == nil
}

func (user *User) Clone() *User {
	return &User{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}
