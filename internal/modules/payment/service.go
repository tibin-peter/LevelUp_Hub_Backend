package payment

import (
	"errors"
	"fmt"
	"time"
)

type Service interface {
	CreateOrder(studentID uint, bookingID uint) (*CreateOrderResult, error)
	VerifyRequest(studentID uint, req VerifyRequest) error
	GetStudentPayment(studentID uint) ([]PaymentSummary, error)
	GetMentorEarnings(mentorID uint) (*MentorEarningsResult, error)
	RequestWithdraw(mentorID uint, amount float64) error
	CreditWallet(userID uint, amount float64, txnType string, refID uint) error
	ReleaseEscrow(bookingID uint) error
	RefundEscrow(bookingID uint) error
	// CleanupStaleBookings()
}

type service struct {
	repo    Repository
	booking BookingPort
	rzp     RazorpayClient
	keyID   string
}

func NewService(r Repository, b BookingPort, rzp RazorpayClient, key string) Service {
	return &service{
		repo:    r,
		booking: b,
		rzp:     rzp,
		keyID:   key,
	}
}

type CreateOrderResult struct {
	OrderID  string
	Amount   int64
	Currency string
	Key      string
}

// create order
func (s *service) CreateOrder(studentID uint, bookingID uint) (*CreateOrderResult, error) {
	booking, err := s.booking.GetBookingByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	if booking.StudentID != studentID {
		return nil, errors.New("unauthorized")
	}
	if booking.Status != "pending_payment" {
		return nil, errors.New("not payable for this one")
	}
	existingPayment, _ := s.repo.GetByBookingID(bookingID)
	if existingPayment != nil && existingPayment.Status == "created" {
		return &CreateOrderResult{
			OrderID:  existingPayment.RazorpayOrderID,
			Amount:   existingPayment.Amount,
			Currency: existingPayment.Currency,
			Key:      s.keyID,
		}, nil
	}
	amount := int64(booking.Price * 100)
	idStr := fmt.Sprintf("booking_%d", booking.ID)
	order, err := s.rzp.CreateOrder(amount, "INR", idStr)
	if err != nil {
		return nil, err
	}
	p := &Payment{
		BookingID:       booking.ID,
		StudentID:       studentID,
		MentorID:        booking.MentorID,
		MentorUserID:    booking.MentorUserID,
		Amount:          amount,
		Currency:        "INR",
		RazorpayOrderID: order.ID,
		Status:          "created",
	}
	if err := s.repo.CreatePayment(p); err != nil {
		return nil, err
	}
	return &CreateOrderResult{
		OrderID:  order.ID,
		Amount:   amount,
		Currency: "INR",
		Key:      s.keyID,
	}, nil
}

func (s *service) VerifyRequest(studentID uint, req VerifyRequest) error {
	ok := s.rzp.VerifySignature(
		req.OrderID,
		req.PaymentID,
		req.Signature,
	)

	if !ok {
		return errors.New("invalid signature")
	}

	payment, err := s.repo.GetPaymentByOrderID(req.OrderID)
	if err != nil || payment == nil {
		return errors.New("payment not found")
	}

	payment.Status = "paid"
	payment.RazorpayPaymentID = string(req.PaymentID)
	payment.PaidAt = time.Now()

	if err := s.repo.UpdatePayment(payment); err != nil {
		return err
	}

	if err := s.booking.MarkBookingPaid(payment.BookingID); err != nil {
		return err
	}

	return nil
}

// student payments
func (s *service) GetStudentPayment(studentID uint) ([]PaymentSummary, error) {
	return s.repo.ListStudentPayments(studentID)
}

type MentorEarningsResult struct {
	Balance float64
	History []WalletTransaction
}

func (s *service) GetMentorEarnings(mentorID uint) (*MentorEarningsResult, error) {

	wallet, err := s.repo.GetWalletByUserID(mentorID)
	if err != nil {
		return nil, err
	}

	var balance float64
	if wallet != nil {
		balance = wallet.Balance
	}

	history, err := s.repo.ListWalletTransactionByUser(mentorID)
	if err != nil {
		return nil, err
	}

	return &MentorEarningsResult{
		Balance: balance,
		History: history,
	}, nil
}

// for withdrawal

func (s *service) RequestWithdraw(mentorID uint, amount float64) error {

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	wallet, err := s.repo.GetWalletByUserID(mentorID)
	if err != nil || wallet == nil {
		return errors.New("wallet not found")
	}

	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	req := &WithdrawRequest{
		MentorID:    mentorID,
		Amount:      amount,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.repo.CreateWithdraw(req); err != nil {
		return err
	}

	wallet.Balance -= amount
	if err := s.repo.UpdateWallet(wallet); err != nil {
		return err
	}

	return nil
}

//// helper for credit balance into mentor and admin

func (s *service) CreditWallet(userID uint, amount float64, txnType string, refID uint) error {

	wallet, err := s.repo.GetWalletByUserID(userID)
	if err != nil {
		return err
	}

	if wallet == nil {
		wallet = &Wallet{
			UserID:  userID,
			Balance: amount,
		}
		if err := s.repo.CreateWallet(wallet); err != nil {
			return err
		}
	} else {
		wallet.Balance += amount
		if err := s.repo.UpdateWallet(wallet); err != nil {
			return err
		}
	}

	txn := &WalletTransaction{
		UserID:      userID,
		Amount:      amount,
		Type:        txnType,
		Source:      "booking",
		ReferenceID: refID,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateWalletTransaction(txn)
}

func (s *service) ReleaseEscrow(bookingID uint) error {

	payment, err := s.repo.GetByBookingID(bookingID)
	if err != nil {
		return err
	}

	if payment.Status != PaymentPaid {
		return errors.New("not in escrow")
	}

	total := float64(payment.Amount) / 100

	adminShare := total * 0.10
	mentorShare := total - adminShare

	// credit mentor
	if err := s.CreditWallet(payment.MentorUserID, mentorShare, "earning", bookingID); err != nil {
		return err
	}

	// credit admin (example userID=1)
	if err := s.CreditWallet(1, adminShare, "commission", bookingID); err != nil {
		return err
	}

	payment.Status = PaymentReleased
	return s.repo.UpdatePayment(payment)
}

func (s *service) RefundEscrow(bookingID uint) error {

	payment, err := s.repo.GetByBookingID(bookingID)
	if err != nil {
		return err
	}

	// if err := s.rzp.Refund(payment.RazorpayPaymentID); err != nil {
	// 	return err
	// }

	// Rule 5: Refund to Student Wallet (Better than Bank Transfer for speed)
	amount := float64(payment.Amount) / 100
	err = s.CreditWallet(payment.StudentID, amount, "refund", bookingID)
	if err != nil {
		return err
	}

	payment.Status = PaymentRefunded
	return s.repo.UpdatePayment(payment)
}

// func (s *service) CleanupStaleBookings() {
// 	// Case 1: Release slots where student started payment but never finished
// 	expiryTime := time.Now().Add(-15 * time.Minute)
// 	staleBookings, _ := s.repo.GetStalePendingBookings(expiryTime)

// 	for _, b := range staleBookings {
// 		fmt.Printf("[CLEANUP] Booking %d timed out. Freeing slot %d\n", b.ID, b.SlotID)
// 		s.booking.UpdateStatus(b.ID, "expired")
// 		s.booking.FreeSlot(b.SlotID)
// 	}

// 	// Case 2: Auto-Refund if Mentor ignored the request and session start time passed
// 	missedBookings, _ := s.repo.GetUnapprovedPastBookings(time.Now())

// 	for _, b := range missedBookings {
// 		fmt.Printf("[CLEANUP] Mentor missed booking %d. Triggering auto-refund.\n", b.ID)

// 		s.booking.UpdateStatus(b.ID, "mentor_missed")

// 		err := s.RefundEscrow(b.ID)
// 		if err != nil {
// 			fmt.Printf("[ERROR] Auto-refund failed for booking %d: %v\n", b.ID, err)
// 			continue
// 		}

// 		s.booking.FreeSlot(b.SlotID)
// 	}
// }
