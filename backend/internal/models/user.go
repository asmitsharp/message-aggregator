package models

import (
	"time"

	"github.com/ashmitsharp/msg-agg/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                string `gorm:"primaryKey"`
	Email             string `gorm:"uniqueKey"`
	PasswordHash      string
	Name              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	MasterUserID      string
	MasterUser        MasterUser
	MatrixUserID      string
	MatrixPassword    string
	MatrixAccessToken string
}

func (u *User) SetPassword(password string) error {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hasedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = utils.GenerateUUID()
	}
	return nil
}
