// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/enum"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
)

func OrderPaid() checkout.OrderPaid {
	faker.SeedGoFake()

	return checkout.NewOrderPaid(
		pkg.OrderID(),
		input.OrderPaidParams{
			PaymentParams: input.PaymentParams{
				AmountPaid:    10*(price.MockPriceStdYear.UnitAmount-50) + 10*(price.MockPriceStdYear.UnitAmount-30) + 5*(price.MockPricePrm.UnitAmount-200) + 5*(price.MockPricePrm.UnitAmount-300),
				ApprovedBy:    gofakeit.Username(),
				Description:   null.StringFrom(gofakeit.Sentence(10)),
				PaymentMethod: pkg.PaymentMethodBank,
				TransactionID: null.StringFrom(pkg.TxnID()),
			},
			Offers: []input.PaymentOfferParams{
				{
					Copies:          10,
					Kind:            enum.OrderKindCreate,
					Price:           price.MockPriceStdYear,
					PriceOffPerCopy: 50,
				},
				{
					Copies:          10,
					Kind:            enum.OrderKindRenew,
					Price:           price.MockPriceStdYear,
					PriceOffPerCopy: 30,
				},
				{
					Copies:          5,
					Kind:            enum.OrderKindCreate,
					Price:           price.MockPricePrm,
					PriceOffPerCopy: 200,
				},
				{
					Copies:          5,
					Kind:            enum.OrderKindRenew,
					Price:           price.MockPricePrm,
					PriceOffPerCopy: 300,
				},
			},
		})
}
