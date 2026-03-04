package payment

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"time"
)

// const (
// 	PaymentCreated  = "created"
// 	PaymentPaid     = "paid"
// 	PaymentReleased = "released"
// 	PaymentRefunded = "refunded"

// 	AdminUserID = 1
// )

type Payment struct {
	ID        uint `gorm:"primaryKey"`
	BookingID uint

	StudentID uint

	MentorProfileID uint
	MentorUserID    uint

	Amount   int64
	Currency string

	RazorpayOrderID   string
	RazorpayPaymentID string

	Status string

	CreatedAt time.Time
	PaidAt    time.Time
}

type PaymentSummary struct {
	ID         uint      `json:"id"`
	Amount     int64     `json:"amount"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status"`
	MentorName string    `json:"mentor_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type Wallet struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint         `gorm:"uniqueIndex" json:"user_id"`
	User   profile.User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletTransaction struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint         `json:"user_id"`
	User   profile.User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
	Source string  `json:"source"`

	ReferenceID uint `json:"reference_id"`

	// Extra fields for frontend display
	Currency   string    `json:"currency" gorm:"-"`
	MentorName string    `json:"mentor_name" gorm:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

type WithdrawRequest struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	MentorID uint `json:"mentor_id"`
	Mentor   profile.User `gorm:"foreignKey:MentorID" json:"mentor,omitempty"`

	Amount float64 `json:"amount"`
	Status string  `json:"status"`

	RequestedAt time.Time `json:"requested_at"`
	ProcessedAt time.Time `json:"processed_at"`
}

// verify payment
type VerifyRequest struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}
type AdminPaymentSummary struct {
	ID              uint      `json:"id"`
	BookingID       uint      `json:"booking_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	MentorName      string    `json:"mentor_name"`
	StudentName     string    `json:"student_name"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	CreatedAt       time.Time `json:"created_at"`
	PaidAt          time.Time `json:"paid_at"`
}

type AdminPaymentOverview struct {
	TotalRevenue    float64 `json:"total_revenue"`
	EscrowHolding   float64 `json:"escrow_holding"`
	TotalRefunded   float64 `json:"total_refunded"`
	TotalReleased   float64 `json:"total_released"`
	PendingWithdraw float64 `json:"pending_withdraw"`
}

type AdminWalletOverview struct {
	CurrentBalance     float64 `json:"current_balance"`
	CommissionEarned   float64 `json:"commission_earned"`
	TotalMentorPayouts float64 `json:"total_mentor_payouts"`
	TotalRefundsGiven  float64 `json:"total_refunds_given"`
}
