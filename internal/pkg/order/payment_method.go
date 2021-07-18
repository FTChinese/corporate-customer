package order

type PaymentMethod string

const (
	PaymentMethodBank   PaymentMethod = "bank"
	PaymentMethodStripe PaymentMethod = "stripe"
	PaymentMethodAli    PaymentMethod = "alipay"
	PaymentMethodWx     PaymentMethod = "wechat"
)
