// +build !production

package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

func MockAssignee() Assignee {
	faker.SeedGoFake()

	return Assignee{
		FtcID:    null.StringFrom(uuid.New().String()),
		UnionID:  null.String{},
		Email:    null.StringFrom(gofakeit.Email()),
		UserName: null.StringFrom(gofakeit.Username()),
	}
}

// MockLicence creates a new licence of this price.
// Deprecated.
func MockLicence(p price.Price) ExpandedLicence {

	return ExpandedLicence{
		Licence: NewLicence(
			p,
			pkg.OrderID(),
			admin.Creator{
				AdminID: "",
				TeamID:  "",
			}),
		Assignee: AssigneeJSON{},
	}
}

func MockInvitation(lic ExpandedLicence) Invitation {
	faker.SeedGoFake()

	inv, err := NewInvitation(input.InvitationParams{
		Email:       gofakeit.Email(),
		Description: null.String{},
		LicenceID:   lic.ID,
	}, admin.PassportClaims{
		AdminID: lic.AdminID,
		TeamID:  null.StringFrom(lic.TeamID),
	})

	if err != nil {
		panic(err)
	}

	return inv
}

type LicBuilder struct {
	price   price.Price
	orderID string
	by      admin.Creator
	to      Assignee
}

// NewLicBuilder creates defaults values for a licence.
func NewLicBuilder(p price.Price) LicBuilder {
	return LicBuilder{
		price:   p,
		orderID: pkg.OrderID(),
		by: admin.Creator{
			AdminID: uuid.NewString(),
			TeamID:  pkg.TeamID(),
		},
		to: Assignee{},
	}
}

// SetOrderID changes the default order id
func (b LicBuilder) SetOrderID(id string) LicBuilder {
	b.orderID = id

	return b
}

// SetCreator changes the default creator
func (b LicBuilder) SetCreator(by admin.Creator) LicBuilder {
	b.by = by
	return b
}

// SetAssignee grant a licence to someone.
func (b LicBuilder) SetAssignee(to Assignee) LicBuilder {
	b.to = to

	return b
}

func (b LicBuilder) Build() ExpandedLicence {
	return ExpandedLicence{
		Licence:  NewLicence(b.price, b.orderID, b.by),
		Assignee: AssigneeJSON{Assignee: b.to},
	}
}
