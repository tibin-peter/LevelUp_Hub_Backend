package courses

import (
	"time"

	"gorm.io/gorm"
)

//table struct for courses
type Course struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(150);not null"`
	ImageURL    string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	Category    string    `gorm:"type:varchar(100);index"`
	Level       string    `gorm:"type:varchar(50)"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

//table struct for mentor relation
type MentorCourse struct {
	ID        uint      `gorm:"primaryKey"`
	MentorID  uint      `gorm:"not null;index;uniqueIndex:idx_mentor_course"`
	CourseID  uint      `gorm:"not null;index;uniqueIndex:idx_mentor_course"`
	CreatedAt time.Time

}