// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/sq"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo() Repo {
	return Repo{
		db: db.MockMySQL().Write,
	}
}

func (r Repo) InsertTeam(t admin.Team) {
	_, err := r.db.NamedExec(admin.StmtCreateTeam, t)
	if err != nil {
		panic(err)
	}
}

func (r Repo) InsertTxnRow(lt checkout.LicenceTransaction) {
	const stmt = `
	INSERT INTO b2b.licence_transaction
	SET transaction_id = :txn_id,
		kind = :kind,
		licence_to_renew = :licence_to_renew,
		order_id = :order_id,
		price_id = :price_id,
		admin_id = :admin_id,
		team_id = :team_id,
		created_utc = :created_utc`

	_, err := r.db.NamedExec(stmt, lt)
	if err != nil {
		log.Fatal(err)
	}
}

func (r Repo) InsertLicence(l licence.Licence) {
	_, err := r.db.NamedExec(checkout.StmtCreateLicence, l)

	if err != nil {
		log.Fatal(err)
	}
}

func (r Repo) InsertInvitation(inv licence.Invitation) {
	_, err := r.db.NamedExec(licence.StmtCreateInvitation, inv)
	if err != nil {
		log.Fatal(err)
	}
}

func (r Repo) InsertMembership(m reader.Membership) {
	_, err := r.db.NamedExec(
		reader.StmtCreateMember,
		m)

	if err != nil {
		log.Fatal(err)
	}
}

func (r Repo) InsertOrderSchema(s checkout.OrderInputSchema) {
	// Insert order row
	_, err := r.db.NamedExec(checkout.StmtCreateOrder, s.OrderRow)
	if err != nil {
		log.Fatal(err)
	}

	// Insert cart items
	for _, v := range s.CartItems {
		_, err := r.db.NamedExec(checkout.StmtInsertCartItem, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create transactions for each copy of licence.
	_, err = r.db.Exec(
		checkout.StmtBulkLicenceTxn(len(s.Transactions)).Build(),
		sq.BuildBulkInsertValues(checkout.BulkLicenceTxn(s.Transactions))...,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func (r Repo) CreateReader(a ReaderAccount) {
	const stmt = `
INSERT INTO cmstmp01.userinfo
SET user_id = :ftc_id,
	wx_union_id = :wx_union_id,
	email = :user_email,
	password = MD5(:password),
	user_name = :user_name,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()
`
	_, err := r.db.NamedExec(stmt, a)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateGrantedLicence creates a licence,
// the user who is granted this licence,
// and the membership derived from this licence.
func (r Repo) CreateGrantedLicence(g GrantedLicence) {
	r.CreateReader(g.Account)
	r.InsertLicence(g.ExpLicence.Licence)
	r.InsertMembership(g.Membership)
}
