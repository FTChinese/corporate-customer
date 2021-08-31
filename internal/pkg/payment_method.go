package pkg

import (
	"database/sql/driver"
	"errors"
)

type PaymentMethod string

const (
	PaymentMethodNULL   PaymentMethod = ""
	PaymentMethodBank   PaymentMethod = "bank"
	PaymentMethodStripe PaymentMethod = "stripe"
	PaymentMethodAli    PaymentMethod = "alipay"
	PaymentMethodWx     PaymentMethod = "wechat"
)

func (p PaymentMethod) Value() (driver.Value, error) {
	if p == "" {
		return nil, nil
	}

	return string(p), nil
}

func (p *PaymentMethod) Scan(src interface{}) error {
	if src == nil {
		*p = PaymentMethodNULL
		return nil
	}

	switch s := src.(type) {
	case []byte:
		*p = PaymentMethod(s)
		return nil

	default:
		return errors.New("incompatible type to scan to PaymentMethod")
	}
}
