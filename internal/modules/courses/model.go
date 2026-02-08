package courses

import "time"

//table struct for courses
type Course struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(150);not null"`
	Description string    `gorm:"type:text"`
	Category    string    `gorm:"type:varchar(100);index"`
	Level       string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//table struct for mentor relation
type MentorCourse struct {
	ID        uint      `gorm:"primaryKey"`
	MentorID  uint      `gorm:"not null;index;uniqueIndex:idx_mentor_course"`
	CourseID  uint      `gorm:"not null;index;uniqueIndex:idx_mentor_course"`
	CreatedAt time.Time

}