package complaints

import "time"

type Complaint struct {
	ID uint `gorm:"primaryKey"`

	UserID uint   `gorm:"index;not null"`
	Role   string `gorm:"not null"`

	Category    string
	Subject     string
	Description string

	Status     string `gorm:"default:'open'"`
	AdminReply string

	CreatedAt time.Time
	UpdatedAt time.Time
}