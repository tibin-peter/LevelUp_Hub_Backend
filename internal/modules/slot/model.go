package slot

import "time"

type MentorSlot struct {
	ID        uint      `gorm:"primaryKey"`
	MentorProfileID  uint      `gorm:"not null;index"`
	StartTime time.Time `gorm:"not null;index"`
	EndTime   time.Time `gorm:"not null;index"`
	IsBooked  bool      `gorm:"default:false"`
	CreatedAt time.Time
}