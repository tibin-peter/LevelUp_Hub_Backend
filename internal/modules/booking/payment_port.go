package booking

type PaymentPort interface {
    ReleaseEscrow(bookingID uint) error
    RefundEscrow(bookingID uint) error
}