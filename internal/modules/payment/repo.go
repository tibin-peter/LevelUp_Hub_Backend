package payment

import (
	"LevelUp_Hub_Backend/internal/repository/generic"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	//payment//
	CreatePayment(p *Payment) error
	UpdatePayment(p *Payment) error
	GetPaymentByOrderID(orderID string) (*Payment, error)
	ListStudentPayments(studentID uint) ([]PaymentSummary, error)
	GetByBookingID(bookingID uint) (*Payment, error)
	SumByMentor(profileID uint) (float64, error)
	SumCumulativeEarnings(profileID uint, result *float64) error

	//wallet//
	GetWalletByUserID(userID uint) (*Wallet, error)
	CreateWallet(w *Wallet) error
	UpdateWallet(w *Wallet) error

	//wallet transation//
	CreateWalletTransaction(t *WalletTransaction) error
	ListWalletTransactionByUser(userID uint) ([]WalletTransaction, error)
	SumWalletTransactionByDate(userID uint, from, to time.Time) (float64, error)

	//withdraw//
	CreateWithdraw(req *WithdrawRequest) error
	ListPendingWithdraws() ([]WithdrawRequest, error)
	ListAllWithdraws() ([]WithdrawRequest, error)
	ListWithdrawalsByMentor(mentorID uint) ([]WithdrawRequest, error)
	GetWithdrawByID(id uint) (*WithdrawRequest, error)
	UpdateWithdraw(req *WithdrawRequest) error

	// Admin Ops
	ListAllPayments(search, status string) ([]AdminPaymentSummary, error)
	GetAdminPaymentOverview() (*AdminPaymentOverview, error)
	GetAdminWalletOverview() (*AdminWalletOverview, error)
}

type repo struct {
	db                    *gorm.DB
	paymentBase           *generic.Repository[Payment]
	walletBase            *generic.Repository[Wallet]
	walletTransactionBase *generic.Repository[WalletTransaction]
	withdrawBase          *generic.Repository[WithdrawRequest]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		db:                    db,
		paymentBase:           generic.NewRepository[Payment](db),
		walletBase:            generic.NewRepository[Wallet](db),
		walletTransactionBase: generic.NewRepository[WalletTransaction](db),
		withdrawBase:          generic.NewRepository[WithdrawRequest](db),
	}
}

//////  payment  //////
func (r *repo) CreatePayment(p *Payment) error {
	return r.paymentBase.Create(p)
}

func (r *repo) UpdatePayment(p *Payment) error {
	return r.paymentBase.Update(p)
}

func (r *repo) GetPaymentByOrderID(orderID string) (*Payment, error) {
	var p Payment
	err := r.db.Where("razorpay_order_id = ?", orderID).First(&p).Error

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) GetByBookingID(id uint) (*Payment, error) {
	var p Payment
	err := r.db.Where("booking_id = ?", id).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &p, err
}

func (r *repo) ListStudentPayments(studentID uint) ([]PaymentSummary, error) {
	var list []PaymentSummary

	err := r.db.Table("payments").
		Select(`
            payments.id,
            payments.amount,
            payments.currency,
            payments.status,
            payments.created_at,
            mentors.name as mentor_name
        `).
		Joins("JOIN mentor_profiles mp ON mp.id = payments.mentor_profile_id").
		Joins("JOIN users mentors ON mentors.id = mp.user_id").
		Where("payments.student_id = ?", studentID).
		Order("payments.created_at DESC").
		Scan(&list).Error

	return list, err
}

/////// Wallet /////////

func (r *repo) GetWalletByUserID(userID uint) (*Wallet, error) {
	var w Wallet
	err := r.db.Where("user_id = ?", userID).First(&w).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &w, err
}

func (r *repo) CreateWallet(w *Wallet) error {
	return r.walletBase.Create(w)
}

func (r *repo) UpdateWallet(w *Wallet) error {
	return r.walletBase.Update(w)
}

////// wallet transaction ///////

func (r *repo) CreateWalletTransaction(t *WalletTransaction) error {
	return r.walletTransactionBase.Create(t)
}

func (r *repo) ListWalletTransactionByUser(userID uint) ([]WalletTransaction, error) {
	var list []WalletTransaction

	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *repo) SumWalletTransactionByDate(userID uint, from, to time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&WalletTransaction{}).Select("COALESCE(SUM(amount),0)").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).Scan(&total).Error

	return total, err
}

////// withdraw ////////

func (r *repo) CreateWithdraw(req *WithdrawRequest) error {
	return r.withdrawBase.Create(req)
}

func (r *repo) ListPendingWithdraws() ([]WithdrawRequest, error) {
	var list []WithdrawRequest
	err := r.db.Preload("Mentor").Where("status = ?", "pending").Order("requested_at ASC").Find(&list).Error
	return list, err
}

func (r *repo) ListAllWithdraws() ([]WithdrawRequest, error) {
	var list []WithdrawRequest
	err := r.db.Preload("Mentor").Order("requested_at DESC").Find(&list).Error
	return list, err
}

func (r *repo) ListWithdrawalsByMentor(mentorID uint) ([]WithdrawRequest, error) {
	var list []WithdrawRequest
	err := r.db.Where("mentor_id = ?", mentorID).Order("requested_at DESC").Find(&list).Error
	return list, err
}

func (r *repo) GetWithdrawByID(id uint) (*WithdrawRequest, error) {
	var w WithdrawRequest
	err := r.db.Preload("Mentor").First(&w, id).Error
	return &w, err
}

func (r *repo) UpdateWithdraw(req *WithdrawRequest) error {
	return r.withdrawBase.Update(req)
}

func (r *repo) SumByMentor(profileID uint) (float64, error) {
	var sum float64
	err := r.db.Model(&Payment{}).
		Where("mentor_profile_id=?", profileID).
		Select("COALESCE(SUM(amount),0)").
		Scan(&sum).Error
	return sum, err
}

func (r *repo) SumCumulativeEarnings(profileID uint, result *float64) error {
	var sum int64
	// Sum both Paid (in escrow) and Released (paid to wallet) payments
	err := r.db.Model(&Payment{}).
		Where("mentor_profile_id = ? AND status IN ?", profileID, []string{"paid", "released"}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	if err != nil {
		return err
	}
	*result = float64(sum) / 100
	return nil
}

// Admin Operations Implementation

func (r *repo) ListAllPayments(search, status string) ([]AdminPaymentSummary, error) {
	var list []AdminPaymentSummary

	query := r.db.Table("payments").
		Select(`
            payments.id,
            payments.booking_id,
            payments.amount,
            payments.currency,
            payments.status,
            payments.created_at,
            payments.paid_at,
            payments.razorpay_order_id,
            mentors.name as mentor_name,
            students.name as student_name
        `).
		Joins("JOIN mentor_profiles mp ON mp.id = payments.mentor_profile_id").
		Joins("JOIN users mentors ON mentors.id = mp.user_id").
		Joins("JOIN users students ON students.id = payments.student_id")

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("mentors.name ILIKE ? OR students.name ILIKE ? OR CAST(payments.id AS TEXT) ILIKE ? OR payments.razorpay_order_id ILIKE ?",
			s, s, s, s)
	}

	if status != "" {
		query = query.Where("payments.status = ?", status)
	}

	err := query.Order("payments.created_at DESC").Scan(&list).Error
	return list, err
}

func (r *repo) GetAdminPaymentOverview() (*AdminPaymentOverview, error) {
	var overview AdminPaymentOverview

	// Total Revenue (all paid/released)
	r.db.Model(&Payment{}).
		Where("status IN ?", []string{PaymentPaid, PaymentReleased}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&overview.TotalRevenue)
	overview.TotalRevenue /= 100

	// Escrow Holding (paid but not released/refunded)
	r.db.Model(&Payment{}).
		Where("status = ?", PaymentPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&overview.EscrowHolding)
	overview.EscrowHolding /= 100

	// Total Refunded
	r.db.Model(&Payment{}).
		Where("status = ?", PaymentRefunded).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&overview.TotalRefunded)
	overview.TotalRefunded /= 100

	// Total Released to Mentors
	r.db.Model(&Payment{}).
		Where("status = ?", PaymentReleased).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&overview.TotalReleased)
	overview.TotalReleased /= 100

	// Pending Withdrawals
	r.db.Model(&WithdrawRequest{}).
		Where("status = ?", "pending").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&overview.PendingWithdraw)

	return &overview, nil
}

func (r *repo) GetAdminWalletOverview() (*AdminWalletOverview, error) {
	var overview AdminWalletOverview

	// Current Balance
	var wallet Wallet
	r.db.Where("user_id = ?", AdminUserID).First(&wallet)
	overview.CurrentBalance = wallet.Balance

	// Commission Earned (Sum of positive transactions from platform fees)
	// Or more accurately, total revenue * 0.1 for completed ones, but let's use transaction history
	r.db.Model(&WalletTransaction{}).
		Where("user_id = ? AND type = ?", AdminUserID, "booking_received").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&overview.CommissionEarned)
	
	// Adjust Commission: total received - (payouts to mentors from those bookings)
	// Simplification: In this system, commission is what stays in Admin Wallet.
	
	// Total Payouts to Mentors
	r.db.Model(&WalletTransaction{}).
		Where("type = ?", "mentor_payout").
		Select("COALESCE(ABS(SUM(amount)), 0)").
		Scan(&overview.TotalMentorPayouts)

	// Total Refunds Given
	r.db.Model(&WalletTransaction{}).
		Where("type = ?", "refund").
		Select("COALESCE(ABS(SUM(amount)), 0)").
		Scan(&overview.TotalRefundsGiven) 

	return &overview, nil
}