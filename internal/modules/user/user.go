package user

import "time"

//user table model
type User struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Role          string `gorm:"not null"`
	IsVerified    bool
	ProfilePicURL string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

//OTP model
type OTP struct {
	ID        uint `gorm:"primaryKey"`
	Email     string
	Code      string
	ExpiresAt time.Time
}
