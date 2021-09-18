// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"log"
)

// NewLicBuilder creates defaults values for a licence.
func (a Admin) NewLicBuilder(p price.Price) LicenceBuilder {
	return LicenceBuilder{
		price:   p,
		txnID:   pkg.TxnID(),
		admin:   a,
		persona: Persona{},
	}
}

func (a Admin) StdLicenceBuilder() LicenceBuilder {
	return a.NewLicBuilder(price.MockPriceStdYear)
}

func (a Admin) PrmLicenceBuilder() LicenceBuilder {
	return a.NewLicBuilder(price.MockPricePrm)
}

type LicenceBuilder struct {
	price   price.Price
	txnID   string
	admin   Admin
	to      licence.Assignee // Deprecated
	persona Persona
}

func (b LicenceBuilder) SetPersona(p Persona) LicenceBuilder {
	b.persona = p
	return b
}

// Build generates mocking licence.
// The licence will have invitation if assignee is set.
func (b LicenceBuilder) Build() licence.Licence {
	lic := licence.
		NewLicence(b.price, b.txnID, b.admin.Creator())

	if b.persona.IsEmpty() {
		return lic
	}

	inv, err := licence.NewInvitation(input.InvitationParams{
		Email:       b.persona.email,
		Description: null.StringFrom(gofakeit.Sentence(10)),
		LicenceID:   lic.ID,
	}, b.admin.PassportClaims())

	if err != nil {
		panic(err)
	}

	return lic.WithGranted(b.persona.Assignee(), inv)
}

func (b LicenceBuilder) BuildExpanded() licence.ExpandedLicence {
	return licence.ExpandedLicence{
		Licence:  b.Build(),
		Assignee: licence.AssigneeJSON{Assignee: b.persona.Assignee()},
	}
}

type GrantedLicence struct {
	Account    ReaderAccount
	ExpLicence licence.ExpandedLicence
	Membership reader.Membership
}

func (b LicenceBuilder) BuildGranted() GrantedLicence {
	if b.persona.IsEmpty() {
		log.Fatal("granted licence must have assignee")
	}

	expLic := b.BuildExpanded()

	return GrantedLicence{
		Account:    b.persona.Account(),
		ExpLicence: expLic,
		Membership: expLic.NewMembership(
			reader.UserIDs{
				CompoundID: b.persona.ftcID,
				FtcID:      null.StringFrom(b.persona.ftcID),
				UnionID:    null.String{},
			},
			addon.AddOn{},
		),
	}
}
