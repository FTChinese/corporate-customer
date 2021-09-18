package txrepo

import (
	"database/sql"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
)

// CreateMember creates a new membership.
func (tx TxRepo) CreateMember(m reader.Membership) error {

	_, err := tx.NamedExec(
		reader.StmtCreateMember,
		m)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMember updates existing membership.
func (tx TxRepo) UpdateMember(m reader.Membership) error {
	_, err := tx.NamedExec(
		reader.StmtUpdateMember,
		m)

	if err != nil {
		return err
	}

	return nil
}

func (tx TxRepo) UpsertMember(m reader.Membership, isNew bool) error {
	if isNew {
		return tx.CreateMember(m)
	}

	return tx.UpdateMember(m)
}

// SaveInvoice inserts a new invoice to db.
// This happens whenever a valid one-time purchase
// membership is changed by a licence.
func (tx TxRepo) SaveInvoice(inv reader.Invoice) error {
	_, err := tx.NamedExec(reader.StmtCreateInvoice, inv)
	if err != nil {
		return err
	}

	return nil
}

// LockMember retrieves a user's membership by a compound id, which might be ftc id or union id.
// Use SQL FIND_IN_SET(compoundId, vip_id, vip) to verify it against two columns.
// Returns zero value of membership if not found.
func (tx TxRepo) LockMember(compoundID string) (reader.Membership, error) {
	var m reader.Membership

	err := tx.Get(
		&m,
		reader.StmtLockMember,
		compoundID,
	)

	if err != nil && err != sql.ErrNoRows {
		return m, err
	}

	// Treat a non-existing member as a valid value.
	return m.Sync(), nil
}
