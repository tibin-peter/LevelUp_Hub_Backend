package favorites

import "time"

type Favorite struct {
	ID uint `gorm:"primaryKey"`

	StudentID       uint `gorm:"not null;uniqueIndex:idx_student_mentor"`
	MentorProfileID uint `gorm:"not null;uniqueIndex:idx_student_mentor"`

	CreatedAt time.Time
}