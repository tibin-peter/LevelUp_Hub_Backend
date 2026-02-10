package booking

import "time"

type Booking struct {
	ID              uint   `gorm:"primaryKey"`
	StudentID       uint   `gorm:"not null;index"`
	MentorProfileID uint   `gorm:"not null;index"`
	SlotID          uint   `gorm:"not null;uniqueIndex"`
	Status          string `gorm:"default;'pending'"`
	Price           float64
	CreatedAt       time.Time
}