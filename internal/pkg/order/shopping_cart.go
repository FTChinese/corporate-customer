package order

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
)

type CarItem struct {
	Price     price.Price       `json:"price"`
	NewCopies int64             `json:"newCopies"`
	Renewals  []licence.Licence `json:"renewals"`
}

type ShoppingCart struct {
	Items       []CarItem `json:"items"`
	ItemCount   int64     `json:"itemCount"`
	TotalAmount int64     `json:"total_amount"`
}
