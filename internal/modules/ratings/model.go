package ratings

import "time"

type Rating struct {
    ID        uint      `gorm:"primaryKey"`
    
    BookingID uint      `gorm:"not null;uniqueIndex"`
    StudentID uint      `gorm:"not null;index"`
    MentorProfileID  uint      `gorm:"not null;index"`

    Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
    Review    string

    CreatedAt time.Time
}