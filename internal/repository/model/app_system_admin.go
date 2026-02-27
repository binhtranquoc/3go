package model

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAdminInactive         = errors.New("tài khoản admin chưa được kích hoạt")
	ErrAdminPasswordMismatch = errors.New("mật khẩu không chính xác")
)

type SystemAdmin struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	Department   string    `json:"department"`
	IsActive     bool      `json:"is_active"`
	LastLoginAt  *string   `json:"last_login_at"`
	BaseModel
}

func (admin *SystemAdmin) SetPassword(rawPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.PasswordHash = string(hash)
	return nil
}

func (admin *SystemAdmin) VerifyPassword(rawPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(rawPassword)) == nil
}

func (admin *SystemAdmin) CanLogin() error {
	if !admin.IsActive {
		return ErrAdminInactive
	}
	return nil
}
