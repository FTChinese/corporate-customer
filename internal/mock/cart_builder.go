// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
)

func (a Admin) CartBuilder() CartBuilder {
	return CartBuilder{
		admin: a,
		store: map[string]checkout.CartItem{},
	}
}

type CartBuilder struct {
	admin Admin
	store map[string]checkout.CartItem
}

func (b CartBuilder) AddNewN(p price.Price, n int) CartBuilder {
	if n == 0 {
		return b
	}

	item, ok := b.store[p.ID]
	if !ok {
		item = checkout.CartItem{
			Price:     p,
			NewCopies: 0,
			Renewals:  nil,
		}
	}
	item.NewCopies = item.NewCopies + int64(n)
	b.store[p.ID] = item

	return b
}

func (b CartBuilder) AddNewStandardN(n int) CartBuilder {
	return b.AddNewN(price.MockPriceStdYear, n)
}

func (b CartBuilder) AddNewPremiumN(n int) CartBuilder {
	return b.AddNewN(price.MockPricePrm, n)
}

func (b CartBuilder) AddRenewal(lic licence.ExpandedLicence) CartBuilder {
	item, ok := b.store[lic.LatestPrice.ID]
	if !ok {
		if !ok {
			item = checkout.CartItem{
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
		item = checkout.CartItem{
			Price:     l.LatestPrice,
			NewCopies: 0,
			Renewals:  nil,
		}
	}
	item.Renewals = append(item.Renewals, lics...)
	b.store[l.LatestPrice.ID] = item

	return b
}

// AddRenewalStandardN adds a licence to renew into cart.
// NOTE the renewal licence added this way are not granted to anyone.
func (b CartBuilder) AddRenewalStandardN(n int) CartBuilder {
	var lic = make([]licence.ExpandedLicence, 0)
	for i := 0; i < n; i++ {
		lic = append(lic, b.admin.StdLicenceBuilder().BuildExpanded())
	}

	return b.AddRenewalN(lic)
}

func (b CartBuilder) AddRenewalPremiumN(n int) CartBuilder {
	var lic = make([]licence.ExpandedLicence, 0)
	for i := 0; i < n; i++ {
		lic = append(lic, b.admin.PrmLicenceBuilder().BuildExpanded())
	}

	return b.AddRenewalN(lic)
}

func (b CartBuilder) Build() checkout.ShoppingCart {
	var items = make([]checkout.CartItem, 0)

	var totalCount int
	var totalAmount float64
	for _, v := range b.store {
		count := int(v.NewCopies) + len(v.Renewals)
		amount := v.Price.UnitAmount * float64(count)

		totalCount += count
		totalAmount += amount

		items = append(items, v)
	}

	return checkout.ShoppingCart{
		Items:       items,
		ItemCount:   int64(totalCount),
		TotalAmount: totalAmount,
	}
}

func (b CartBuilder) BuildOrderSchema() checkout.OrderInputSchema {
	return checkout.
		NewOrderSchemaBuilder(b.Build(), b.admin.PassportClaims()).
		Build()
}
