// +build !production

package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/price"
)

type MockRepo struct {
	dbs db.ReadWriteMyDBs
}

func MockNewRepo() MockRepo {
	return MockRepo{
		dbs: db.MockMySQL(),
	}
}

func (r MockRepo) MustCreateLicence(l licence.Licence) {

	_, err := r.dbs.Write.NamedExec(licence.StmtCreateLicence, l)
	if err != nil {
		panic(err)
	}
}

func (r MockRepo) MustCreateInvitation(inv licence.Invitation) {
	_, err := r.dbs.Write.NamedExec(licence.StmtCreateInvitation, inv)
	if err != nil {
		panic(err)
	}
}

func (r MockRepo) MustCreateMember(m reader.Membership) {
	_, err := r.dbs.Write.NamedExec(reader.StmtCreateMember, m)
	if err != nil {
		panic(err)
	}
}

func (r MockRepo) MustCreateInvitedLicence(a licence.Assignee) licence.ExpandedLicence {
	lic := licence.MockLicence(price.MockPriceStdYear)
	inv := licence.MockInvitation(lic)
	baseLic := lic.WithInvitation(inv)

	r.MustCreateLicence(baseLic)
	r.MustCreateInvitation(inv)

	return licence.ExpandedLicence{
		Licence:  baseLic,
		Assignee: licence.AssigneeJSON{Assignee: a},
	}
}

func (r MockRepo) MustCreateGrantedLicence(a licence.Assignee) licence.ExpandedLicence {
	lic := licence.MockLicence(price.MockPriceStdYear)
	inv := licence.MockInvitation(lic).Accepted()
	baseLic := lic.Granted(a, inv)

	result := licence.NewGrantResult(
		licence.ExpandedLicence{
			Licence:  baseLic,
			Assignee: licence.AssigneeJSON{Assignee: a},
		},
		reader.Membership{},
	)

	r.MustCreateLicence(result.Licence.Licence)
	r.MustCreateInvitation(result.Licence.LatestInvitation.Invitation)
	r.MustCreateMember(result.Membership)

	return licence.ExpandedLicence{
		Licence:  baseLic,
		Assignee: licence.AssigneeJSON{Assignee: a},
	}
}
