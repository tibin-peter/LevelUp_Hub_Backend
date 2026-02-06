package profile

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

//mentor details table
type MentorProfile struct{
	ID    uint `gorm:"primaryKey"`
	UserID  uint `gorm:"uniqueIndex;not null"`
	Bio string `gorm:"type:text"`
	Skills string `gorm:"type:text"`
	Languages string `gorm:"type:text"`
	HourlyPrice     float64 `gorm:"not null"`
	ExperienceYears int     `gorm:"not null"`
	RatingAvg    float64 `gorm:"default:0"`
	TotalReviews int     `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
