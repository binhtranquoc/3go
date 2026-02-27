package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"password_hash"`
	BaseModel
}

func (account *Account) SetPassword(rawPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.PasswordHash = string(hash)
	return nil
}

func (account *Account) VerifyPassword(rawPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(rawPassword)) == nil
}
