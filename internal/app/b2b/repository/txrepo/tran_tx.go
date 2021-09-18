package txrepo

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/go-rest/enum"
)

// LockLicenceTxn retrieves a row from licence_transaction.
func (tx TxRepo) LockLicenceTxn(id string) (checkout.LicenceTransaction, error) {
	var lq checkout.LicenceTransaction
	err := tx.Get(&lq, checkout.StmtLockLicenceTxn, id)
	if err != nil {
		return checkout.LicenceTransaction{}, err
	}

	return lq, nil
}

// FinalizeLicenceTxn updates a row in licence queue table
// after licence is creates/renewed.
func (tx TxRepo) FinalizeLicenceTxn(lt checkout.LicenceTransaction) error {
	_, err := tx.NamedExec(
		checkout.StmtFinalizeLicenceTxn,
		lt)

	if err != nil {
		return err
	}

	return nil
}

// LockLicenceCMS locks and retrieves a row from licence table
// before updating it in CMS.
// If the passed-in id is an empty string, this method returns zero
// value of ExpandedLicence. This is used when creating a new licence
// since a new one does not have existing id.
func (tx TxRepo) LockLicenceCMS(id string) (licence.Licence, error) {
	if id == "" {
		return licence.Licence{}, nil
	}

	var l licence.Licence
	err := tx.Get(
		&l,
		licence.StmtLockLicenceCMS,
		id)
	if err != nil {
		return licence.Licence{}, err
	}

	return l, nil
}

// CreateLicence inserts a new row into licence table.
func (tx TxRepo) CreateLicence(lic licence.Licence) error {
	_, err := tx.NamedExec(checkout.StmtCreateLicence, lic)
	if err != nil {
		return err
	}

	return nil
}

// RenewLicence extends a licence and may optionally revoke
// it if it is granted to someone, expired and the assignee
// switched to auto-renewal subscription.
func (tx TxRepo) RenewLicence(lic licence.Licence) error {
	_, err := tx.NamedExec(checkout.StmtRenewLicence, lic)
	if err != nil {
		return err
	}

	return nil
}

// UpsertLicence insert or update a licence based on
// what kind of licence is being created.
func (tx TxRepo) UpsertLicence(k enum.OrderKind, lic licence.Licence) error {
	switch k {
	case enum.OrderKindCreate:
		return tx.CreateLicence(lic)

	case enum.OrderKindRenew:
		return tx.RenewLicence(lic)
	}

	return errors.New("licence upsert only support create or renew kind")
}
