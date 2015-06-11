package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Password string `json:"-"`
}

func (u *User) SetPassword(p string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(b)
	return nil
}

func (u *User) CheckPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return false
	}
	return true
}
