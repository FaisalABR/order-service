package constants

type OrderStatus int
type OrderStatusString string

const (
	Pending        OrderStatus = 100
	PendingPayment OrderStatus = 200
	PaymentSuccess OrderStatus = 300
	Expired        OrderStatus = 400

	PendingString        OrderStatusString = "pending"
	PendingPaymentString OrderStatusString = "pending-payment"
	PaymentSuccessString OrderStatusString = "payment-success"
	ExpiredString        OrderStatusString = "Expired"
)

var mapOrderStatusStringToInt = map[OrderStatusString]OrderStatus{
	PendingString:        Pending,
	PendingPaymentString: PendingPayment,
	PaymentSuccessString: PaymentSuccess,
	ExpiredString:        Expired,
}

var mapOrderStatusIntToString = map[OrderStatus]OrderStatusString{
	Pending:        PendingString,
	PendingPayment: PendingPaymentString,
	PaymentSuccess: PaymentSuccessString,
	Expired:        ExpiredString,
}

func (p OrderStatus) GetStatusString() OrderStatusString {
	return mapOrderStatusIntToString[p]
}

func (p OrderStatusString) GetStatus() OrderStatus {
	return mapOrderStatusStringToInt[p]
}
