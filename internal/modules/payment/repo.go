package payment

import (
	"LevelUp_Hub_Backend/internal/repository/generic"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	//payment//
	CreatePayment(p *Payment)error
	UpdatePayment(p *Payment)error
	GetPaymentByOrderID(orderID string)(*Payment,error)
	ListStudentPayments(studentID uint)([]PaymentSummary,error)
	GetByBookingID(bookingID uint) (*Payment, error)
	SumByMentor(profileID uint) (float64, error)

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
	UpdateWithdraw(req *WithdrawRequest) error

}

type repo struct {
	db   *gorm.DB
	paymentBase *generic.Repository[Payment]
  walletBase *generic.Repository[Wallet]
  walletTransactionBase *generic.Repository[WalletTransaction]
  withdrawBase *generic.Repository[WithdrawRequest]
}


func NewRepository(db *gorm.DB)Repository{
	return &repo{
		db: db,
		paymentBase: generic.NewRepository[Payment](db),
		walletBase:   generic.NewRepository[Wallet](db),
		walletTransactionBase:      generic.NewRepository[WalletTransaction](db),
		withdrawBase: generic.NewRepository[WithdrawRequest](db),
	}
}


//////  payment  //////
func(r *repo)CreatePayment(p *Payment)error{
	return r.paymentBase.Create(p)
}

func(r *repo)UpdatePayment(p *Payment)error{
	return r.paymentBase.Update(p)
}

func(r *repo)GetPaymentByOrderID(orderID string)(*Payment,error){
	var p Payment
	err:=r.db.Where("razorpay_order_id = ?",orderID).First(&p).Error
	
	if err!=nil{
		return nil,err
	}
	return &p,nil
}

func (r *repo) GetByBookingID(id uint) (*Payment, error) {
    var p Payment
    err := r.db.Where("booking_id = ?", id).First(&p).Error
    return &p, err
}

func (r *repo) ListStudentPayments(studentID uint) ([]PaymentSummary, error) {
    var list []PaymentSummary

    err := r.db.Table("payments").
        Select("payments.id, payments.amount, payments.currency, payments.status, payments.created_at, mentors.name as mentor_name").
        Joins("LEFT JOIN users AS mentors ON mentors.id = payments.mentor_id").
        Where("payments.student_id = ?", studentID).
        Order("payments.created_at DESC").
        Scan(&list).Error

    return list, err
}

/////// Wallet /////////

func(r *repo)GetWalletByUserID(userID uint)(*Wallet,error){
	var w Wallet
	err:=r.db.Where("user_id = ?",userID).First(&w).Error
	if err == gorm.ErrRecordNotFound{
		return nil,nil
	}
	return &w,err
}

func(r *repo)CreateWallet(w *Wallet)error{
	return r.walletBase.Create(w)
}

func (r *repo)UpdateWallet(w *Wallet)error{
	return r.walletBase.Update(w)
}


////// wallet transaction ///////

func (r *repo) CreateWalletTransaction(t *WalletTransaction) error {
	return r.walletTransactionBase.Create(t)
}

func (r *repo)ListWalletTransactionByUser(userID uint)([]WalletTransaction,error){
	var list []WalletTransaction

	err:=r.db.Where("user_id = ?",userID).Order("created_at DESC").Find(&list).Error
	return list,err
}

func (r *repo)SumWalletTransactionByDate(userID uint,from,to time.Time)(float64,error){
	var total float64
	err:=r.db.Model(&WalletTransaction{}).Select("COALESCE(SUM(amount),0)").
	Where("user_id = ? AND created_at BETWEEN ? AND ?",userID,from,to).Scan(&total).Error

	return total,err
}

////// withdraw ////////

func (r *repo) CreateWithdraw(req *WithdrawRequest) error {
	return r.withdrawBase.Create(req)
}

func (r *repo) ListPendingWithdraws() ([]WithdrawRequest, error) {
	var list []WithdrawRequest

	err := r.db.
		Where("status = ?", "pending").
		Order("requested_at ASC").
		Find(&list).Error

	return list, err
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