package payment

const (
    PaymentCreated  = "created"
    PaymentPaid     = "paid"      // money captured (escrow)
    PaymentReleased = "released"  // paid to mentor
    PaymentRefunded = "refunded"
)