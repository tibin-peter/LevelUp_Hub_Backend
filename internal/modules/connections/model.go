package connections

import "time"

type Connection struct {
	ID uint `gorm:"primaryKey"`

	StudentID       uint `gorm:"not null;uniqueIndex:idx_conn"`
	MentorProfileID uint `gorm:"not null;uniqueIndex:idx_conn"`

	CreatedAt time.Time
}