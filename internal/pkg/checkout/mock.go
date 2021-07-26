// +build !production

package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

func MockCartItem(tier enum.Tier) CartItem {
	var p price.Price
	switch tier {
	case enum.TierStandard:
		p = price.MockPriceStdYear
	case enum.TierPremium:
		p = price.MockPricePrm
	}
	newCopies := rand.IntRange(1, 5)
	renewCopies := rand.IntRange(1, 5)

	var renewals = make([]licence.Licence, 0)
	for i := 0; i < renewCopies; i++ {
		renewals = append(renewals, licence.MockLicence(p))
	}

	return CartItem{
		Price:     p,
		NewCopies: int64(newCopies),
		Renewals:  renewals,
	}
}

func MockShoppingCart() ShoppingCart {
	itemStd := MockCartItem(enum.TierStandard)
	stdCopies := itemStd.NewCopies + int64(len(itemStd.Renewals))

	itemPrm := MockCartItem(enum.TierPremium)
	prmCopies := itemPrm.NewCopies + int64(len(itemPrm.Renewals))

	return ShoppingCart{
		Items: []CartItem{
			itemStd,
			itemPrm,
		},
		ItemCount:   stdCopies + prmCopies,
		TotalAmount: float64(stdCopies)*itemStd.Price.UnitAmount + float64(prmCopies)*itemPrm.Price.UnitAmount,
	}
}

func MockOrderInputSchema() OrderInputSchema {
	cart := MockShoppingCart()
	pp := admin.PassportClaims{
		AdminID: uuid.New().String(),
		TeamID:  null.StringFrom(pkg.TeamID()),
	}

	return NewOrderInputSchema(cart, pp)
}
