package txrepo

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/go-rest/enum"
)

// LockLicenceQueue retrieves a row from licence_queue.
func (tx TxRepo) LockLicenceQueue(id string) (checkout.LicenceQueue, error) {
	var lq checkout.LicenceQueue
	err := tx.Get(&lq, checkout.StmtLockLicenceQueue, id)
	if err != nil {
		return checkout.LicenceQueue{}, err
	}

	return lq, nil
}

// FinalizeLicenceQueue updates a row in licence queue table
// after licence is creates/renewed.
func (tx TxRepo) FinalizeLicenceQueue(q checkout.LicenceQueue) error {
	_, err := tx.NamedExec(
		checkout.StmtFinalizeLicenceQueue,
		q)

	if err != nil {
		return err
	}

	return nil
}

// LockBaseLicenceCMS locks and retrieves a row from licence table
// before updating it in CMS.
// If the passed-in id is an empty string, this method returns zero
// value of ExpandedLicence. This is used when creating a new licence
// since a new one does not have existing id.
func (tx TxRepo) LockBaseLicenceCMS(id string) (licence.Licence, error) {
	if id == "" {
		return licence.Licence{}, nil
	}

	var l licence.Licence
	err := tx.Get(
		&l,
		licence.StmtLockBaseLicenceCMS,
		id)
	if err != nil {
		return licence.Licence{}, err
	}

	return l, nil
}

// CreateLicence inserts a new row into licence table.
func (tx TxRepo) CreateLicence(lic licence.Licence) error {
	_, err := tx.NamedExec(licence.StmtCreateLicence, lic)
	if err != nil {
		return err
	}

	return nil
}

func (tx TxRepo) RenewLicence(lic licence.Licence) error {
	_, err := tx.NamedExec(licence.StmtRenewLicence, lic)
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

// RenewMembership updates membership to the latest state
// based the licence linked to it.
func (tx TxRepo) RenewMembership(lic licence.Licence) (checkout.MembershipRenewed, error) {
	if !lic.IsGranted() {
		return checkout.MembershipRenewed{}, nil
	}

	currMmb, err := tx.RetrieveMember(lic.AssigneeID.String)
	if err != nil {
		return checkout.MembershipRenewed{}, err
	}

	// TODO: should we record this error to db?
	if !lic.IsGrantedTo(currMmb) {
		return checkout.MembershipRenewed{}, nil
	}

	newMmb := licence.NewMembership(
		currMmb.UserIDs,
		lic,
		currMmb.AddOn)

	err = tx.UpdateMember(newMmb)
	if err != nil {
		return checkout.MembershipRenewed{}, err
	}

	return checkout.MembershipRenewed{
		Latest:  newMmb,
		Archive: currMmb.Archive(reader.B2BArchiver(reader.ArchiveActionRenew)),
	}, nil
}
