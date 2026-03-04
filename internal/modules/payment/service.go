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
	GetMentorWithdrawals(mentorID uint) ([]WithdrawRequest, error)

	// Admin
	GetAdminPaymentLedger(search, status string) ([]AdminPaymentSummary, error)
	GetAdminPaymentOverview() (*AdminPaymentOverview, error)
	GetAdminWalletOverview() (*AdminWalletOverview, error)
	GetAdminWalletTransactions() ([]WalletTransaction, error)
	GetAdminWithdrawals() ([]WithdrawRequest, error)
	ApproveWithdrawal(id uint) error
	RejectWithdrawal(id uint) error
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
		MentorProfileID: booking.MentorID,
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

	if payment.Status == PaymentPaid {
		return nil
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
	Balance float64             `json:"balance"`
	History []WalletTransaction `json:"history"`
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

	// Populate Currency and MentorName for each transaction from the associated booking payment
	for i := range history {
		if history[i].Source == "booking" && history[i].ReferenceID > 0 {
			p, err := s.repo.GetByBookingID(history[i].ReferenceID)
			if err == nil && p != nil {
				history[i].Currency = p.Currency
				history[i].MentorName = "Session Earning"
			}
		} else if history[i].Type == "mentor_payout" || history[i].Type == "withdraw" {
			history[i].Currency = "INR"
			history[i].MentorName = "Withdrawal"
		} else {
			history[i].Currency = "INR"
		}
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

func (s *service) GetMentorWithdrawals(mentorID uint) ([]WithdrawRequest, error) {
	return s.repo.ListWithdrawalsByMentor(mentorID)
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

	if payment.Status == PaymentReleased {
		return nil // Idempotent
	}

	if payment.Status != PaymentPaid {
		return errors.New("not in escrow or already processed")
	}

	total := float64(payment.Amount) / 100
	adminShare := total * 0.10
	mentorShare := total - adminShare

	// 1. Credit Admin Wallet with Full Amount
	if err := s.CreditWallet(AdminUserID, total, "booking_received", bookingID); err != nil {
		return err
	}

	// Credit Mentor
	if err := s.CreditWallet(payment.MentorUserID, mentorShare, "session_earning", bookingID); err != nil {
		return err
	}

	// Record the deduction in Admin Wallet (Debit)
	if err := s.CreditWallet(AdminUserID, -mentorShare, "mentor_payout", bookingID); err != nil {
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

	if err := s.rzp.Refund(payment.RazorpayPaymentID); err != nil {
		return err
	}

	amount := float64(payment.Amount) / 100
	err = s.CreditWallet(payment.StudentID, amount, "refund", bookingID)
	if err != nil {
		return err
	}

	payment.Status = PaymentRefunded
	return s.repo.UpdatePayment(payment)
}

// Admin Operations

func (s *service) GetAdminPaymentLedger(search, status string) ([]AdminPaymentSummary, error) {
	return s.repo.ListAllPayments(search, status)
}

func (s *service) GetAdminPaymentOverview() (*AdminPaymentOverview, error) {
	return s.repo.GetAdminPaymentOverview()
}

func (s *service) GetAdminWalletOverview() (*AdminWalletOverview, error) {
	return s.repo.GetAdminWalletOverview()
}

func (s *service) GetAdminWalletTransactions() ([]WalletTransaction, error) {
	return s.repo.ListWalletTransactionByUser(AdminUserID)
}

func (s *service) GetAdminWithdrawals() ([]WithdrawRequest, error) {
	return s.repo.ListAllWithdraws()
}

func (s *service) ApproveWithdrawal(id uint) error {
	req, err := s.repo.GetWithdrawByID(id)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return errors.New("request already processed")
	}

	req.Status = "approved"
	req.ProcessedAt = time.Now()
	return s.repo.UpdateWithdraw(req)
}

func (s *service) RejectWithdrawal(id uint) error {
	req, err := s.repo.GetWithdrawByID(id)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return errors.New("request already processed")
	}

	// Refund Mentor Wallet since we deducted it during RequestWithdraw
	wallet, err := s.repo.GetWalletByUserID(req.MentorID)
	if err != nil {
		return err
	}
	wallet.Balance += req.Amount
	if err := s.repo.UpdateWallet(wallet); err != nil {
		return err
	}

	req.Status = "rejected"
	req.ProcessedAt = time.Now()
	return s.repo.UpdateWithdraw(req)
}
