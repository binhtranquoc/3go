package model

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserProfileAlreadyActive = errors.New("tài khoản đã được kích hoạt")
	ErrUserProfileInactive      = errors.New("tài khoản chưa được kích hoạt")
)

type UserProfile struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	FullName  string
	AvatarURL string
	IsActive  bool
	Metadata  []byte
	BaseModel
}

func (userProfile *UserProfile) CanActivate() error {
	if userProfile.IsActive {
		return ErrUserProfileAlreadyActive
	}
	return nil
}

func (userProfile *UserProfile) Activate() {
	userProfile.IsActive = true
}

func (userProfile *UserProfile) CanLogin() error {
	if !userProfile.IsActive {
		return ErrUserProfileInactive
	}
	return nil
}
