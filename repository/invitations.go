package repository

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/guregu/null"
)

var (
	ErrLicenceUnavailable = errors.New("licence is unavailable to grant")
	ErrAlreadyMember      = errors.New("the invitee already has a valid membership")
)

// CreateInvitation creates a new invitation for a licence.
func (env Env) CreateInvitation(input admin.InvitationInput) (admin.Assignee, error) {
	tx, err := env.beginInvTx()
	if err != nil {
		return admin.Assignee{}, err
	}

	// Retrieve the licence. It does not include the Assignee fields
	rawLicence, err := tx.RetrieveLicence(input)
	// There is an not found error here.
	if err != nil {
		_ = tx.Rollback()
		return admin.Assignee{}, err
	}

	// If this licence is not available to grant.
	if !rawLicence.IsAvailable() {
		_ = tx.Rollback()
		return admin.Assignee{}, ErrLicenceUnavailable
	}

	// Try to find the reader account by email.
	// Not found should not be considered an error here.
	invitee, err := env.FindReader(input.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return admin.Assignee{}, err
	}

	// If this reader has a valid membership, disallow
	// granting a new licence.
	if !invitee.Membership.IsExpired() {
		_ = tx.Rollback()
		return admin.Assignee{}, ErrAlreadyMember
	}

	if invitee.FtcID.IsZero() {
		invitee.Email = null.StringFrom(input.Email)
	}

	// Create Invitation instance based on the input data.
	inv, err := input.NewInvitation()
	if err != nil {
		_ = tx.Rollback()
		return admin.Assignee{}, err
	}

	// Update licence with by setting last_invitation column.
	updatedRawLicence, err := rawLicence.WithInvitationSent(inv)
	err = tx.UpdateLicence(updatedRawLicence)
	if err != nil {
		_ = tx.Rollback()
		return admin.Assignee{}, err
	}

	// Save the invitation
	err = tx.SaveInvitation(inv)
	if err != nil {
		return admin.Assignee{}, err
	}

	if err := tx.Commit(); err != nil {
		return admin.Assignee{}, err
	}

	return invitee.Assignee, nil
}

func (env Env) RevokeInvitation(id, teamID string) error {
	// Retrieve the invitation

	// Update revoked field

	// retrieve the licence associated to this invitation

	// Update
	_, err := env.db.Exec(stmt.RevokeInvitation, id, teamID)

	if err != nil {
		return err
	}

	return nil
}

// List invitations shows a list of invitations for a team.
func (env Env) ListInvitations(teamID string) ([]admin.Invitation, error) {
	var invs = make([]admin.Invitation, 0)

	err := env.db.Select(&invs, stmt.ListInvitation, teamID)

	if err != nil {
		return nil, err
	}

	return invs, nil
}

func (env Env) CountInvitation(teamID string) (int64, error) {
	var total int64

	err := env.db.Select(&total, stmt.CountInvitation, teamID)

	if err != nil {
		return total, err
	}

	return total, nil
}
