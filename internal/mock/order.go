// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
)

// ShoppingCart builds a mocking shopping cart for the admin.
func (a Admin) ShoppingCart(items ...checkout.OrderItem) checkout.ShoppingCart {
	b := checkout.NewCartBuilder()
	for _, v := range items {
		b.AddN(v.Price, v.NewCopies)

		for i := 0; i < v.RenewalCopies; i++ {
			b.AddRenewal(a.ExistingExpLicence(v.Price, licence.Assignee{}))
		}
	}

	return b.Build()
}

func (a Admin) OrderByItems(items ...checkout.OrderItem) checkout.Order {
	return a.Order(a.ShoppingCart(items...))
}

func (a Admin) Order(cart checkout.ShoppingCart) checkout.Order {
	return checkout.NewOrder(cart, a.PassportClaims())
}

func (a Admin) OrderSchema(cart checkout.ShoppingCart) checkout.OrderInputSchema {
	return checkout.NewOrderInputSchema(cart, a.PassportClaims())
}
