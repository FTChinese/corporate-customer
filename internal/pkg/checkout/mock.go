// +build !production

package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
)

type CartBuilder struct {
	store map[string]CartItem
}

func NewCartBuilder() CartBuilder {
	return CartBuilder{
		store: map[string]CartItem{},
	}
}

func (b CartBuilder) Add(p price.Price) CartBuilder {
	item, ok := b.store[p.ID]

	if !ok {
		item = CartItem{
			Price:     p,
			NewCopies: 0,
			Renewals:  nil,
		}
	}

	item.NewCopies = item.NewCopies + 1

	b.store[p.ID] = item

	return b
}

func (b CartBuilder) AddN(p price.Price, n int) CartBuilder {
	if n == 0 {
		return b
	}

	item, ok := b.store[p.ID]
	if !ok {
		item = CartItem{
			Price:     p,
			NewCopies: 0,
			Renewals:  nil,
		}
	}
	item.NewCopies = item.NewCopies + int64(n)
	b.store[p.ID] = item

	return b
}

func (b CartBuilder) AddStandardN(n int) CartBuilder {
	return b.AddN(price.MockPriceStdYear, n)
}

func (b CartBuilder) AddPremiumN(n int) CartBuilder {
	return b.AddN(price.MockPricePrm, n)
}

func (b CartBuilder) AddRenewal(lic licence.ExpandedLicence) CartBuilder {
	item, ok := b.store[lic.LatestPrice.ID]
	if !ok {
		if !ok {
			item = CartItem{
				Price:     lic.LatestPrice,
				NewCopies: 0,
				Renewals:  nil,
			}
		}
	}
	item.Renewals = append(item.Renewals, lic)
	b.store[lic.LatestPrice.ID] = item

	return b
}

func (b CartBuilder) AddRenewalN(lics []licence.ExpandedLicence) CartBuilder {
	if lics == nil || len(lics) == 0 {
		return b
	}

	l := lics[0]
	item, ok := b.store[l.LatestPrice.ID]
	if !ok {
		item = CartItem{
			Price:     l.LatestPrice,
			NewCopies: 0,
			Renewals:  nil,
		}
	}
	item.Renewals = append(item.Renewals, lics...)
	b.store[l.LatestPrice.ID] = item

	return b
}

func (b CartBuilder) Build() ShoppingCart {
	var items = make([]CartItem, 0)

	var totalCount int
	var totalAmount float64
	for _, v := range b.store {
		count := int(v.NewCopies) + len(v.Renewals)
		amount := v.Price.UnitAmount * float64(count)

		totalCount += count
		totalAmount += amount

		items = append(items, v)
	}

	return ShoppingCart{
		Items:       items,
		ItemCount:   int64(totalCount),
		TotalAmount: totalAmount,
	}
}
