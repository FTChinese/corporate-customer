package checkout

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/enum"
)

type LicenceBuilder struct {
	Kind           enum.OrderKind
	CurrentLicence licence.Licence
	Price          price.Price
	Order          Order
}

func (b LicenceBuilder) Build() (licence.Licence, error) {
	switch b.Kind {
	case enum.OrderKindCreate:
		return licence.NewLicence(b.Price, b.Order.ID, b.Order.Creator), nil

	case enum.OrderKindRenew:
		return b.CurrentLicence.Renewed(b.Price, b.Order.ID), nil
	}

	return licence.Licence{}, errors.New("unknown order kind")
}
