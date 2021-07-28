// +build !production

package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/db"
)

const mockStmtCreateLicence = `
INSERT INTO b2b.licence
SET id = :licence_id,
	tier = :tier,
	cycle = :cycle,
	current_status = :lic_status,
	creator_id = :creator_id,
	team_id = :team_id,
	current_period_start_utc = :current_period_start_utc,
	current_period_end_utc = :current_period_end_utc,
	start_date_utc = :start_date_utc,
	trial_start_utc = :trial_start_utc,
	trial_end_utc = :trial_end_utc,
	latest_order_id = :latest_order_id,
	latest_price = :latest_price,
	latest_invitation = :latest_invitation,
	assignee_id = :assignee_id,
	created_utc = :created_utc,
	updated_utc = :updated_utc`

type MockRepo struct {
	dbs db.ReadWriteMyDBs
}

func MockNewRepo() MockRepo {
	return MockRepo{
		dbs: db.MockMySQL(),
	}
}

func (r MockRepo) CreateLicence(l licence.Licence) {

	_, err := r.dbs.Write.NamedExec(mockStmtCreateLicence, l)
	if err != nil {
		panic(err)
	}
}

func (r MockRepo) CreateInvitation(inv licence.Invitation) {
	_, err := r.dbs.Write.NamedExec(licence.StmtCreateInvitation, inv)
	if err != nil {
		panic(err)
	}
}

func (r MockRepo) CreateMember(m reader.Membership) {
	_, err := r.dbs.Write.NamedExec(reader.StmtCreateMember, m)
	if err != nil {
		panic(err)
	}
}
