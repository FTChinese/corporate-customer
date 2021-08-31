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
