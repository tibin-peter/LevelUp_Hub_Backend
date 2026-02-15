package payment

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"time"
)

type Payment struct {
	ID        uint `gorm:"primaryKey"`
	BookingID uint

	StudentID uint

	MentorID   uint

	Amount   int64
	Currency string

	RazorpayOrderID   string
	RazorpayPaymentID string

	Status string

	CreatedAt time.Time
	PaidAt    time.Time
}

type PaymentSummary struct {
	ID         uint
	Amount     int64
	Currency   string
	Status     string
	MentorName string
	CreatedAt  time.Time
}

type Wallet struct {
	ID uint `gorm:"primaryKey"`

	UserID uint         `gorm:"uniqueIndex"`
	User   profile.User `gorm:"foreignKey:UserID"`

	Balance   float64
	UpdatedAt time.Time
}

type WalletTransaction struct {
	ID uint `gorm:"primaryKey"`

	UserID uint
	User   profile.User `gorm:"foreignKey:UserID"`

	Amount float64
	Type   string
	Source string

	ReferenceID uint

	CreatedAt time.Time
}

type WithdrawRequest struct {
	ID       uint `gorm:"primaryKey"`
	MentorID uint
	Mentor   profile.User `gorm:"foreignKey:MentorID"`

	Amount float64
	Status string

	RequestedAt time.Time
	ProcessedAt time.Time
}
