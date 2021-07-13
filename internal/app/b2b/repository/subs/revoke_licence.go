package subs

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/model"
)

// RevokeLicence revokes a licence granted to a reader.
// If Licence.Status is not LicStatusGranted, no-ops will
// be performed.
// For a licence waiting for its invitation accepted,
// use RevokeInvitation instead of this one.
// Deprecated
func (env Env) RevokeLicence(id, teamID string) error {
	tx, err := env.beginInvTx()
	if err != nil {
		return err
	}

	lic, err := tx.RetrieveLicence(id, teamID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// If the lic is not granted yet.
	if lic.Status != licence.LicStatusGranted {
		_ = tx.Rollback()
		return nil
	}

	// Get assignee current membership
	mmb, err := tx.RetrieveMembership(lic.Assignee.FtcID.String)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Back up current membership.
	go func() {
		_ = env.TakeSnapshot(model.NewMemberSnapshot(mmb))
	}()

	// Nullify membership's fields.
	newMmb := mmb.WithLicenceRevoked()
	// Update membership.
	err = tx.UpdateMembership(newMmb)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Update lic data.
	revoked := lic.WithInvitationRevoked()
	err = tx.RevokeLicence(revoked.BaseLicence)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
