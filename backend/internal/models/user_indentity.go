package models

import "time"

type UserIdentity struct {
	ID           string `gorm:"primaryKey"`
	MasterUserID string
	Platform     string
	Identifier   string
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MasterUser struct {
	ID         string `gorm:"primaryKey"`
	Name       string
	Email      string
	Phone      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Identities []UserIdentity
}
