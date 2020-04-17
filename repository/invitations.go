package repository

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/guregu/null"
)

var (
	ErrLicenceUnavailable = errors.New("licence is unavailable to grant")
	ErrAlreadyMember      = errors.New("the invitee already has a valid membership")
	ErrInviteeMismatch    = errors.New("an invitation for this licence is already sent another user")
)

// CreateInvitation creates a new invitation for a licence.
// To create an invitation letter, we need the following
// information:
// * Assignee
// * Invitation.Token
// * Plan
func (env Env) CreateInvitation(input admin.InvitationInput) (admin.InvitedLicence, error) {
	tx, err := env.beginInvTx()
	if err != nil {
		return admin.InvitedLicence{}, err
	}

	// Retrieve the licence. It does not include the Assignee fields
	licence, err := tx.RetrieveLicence(input.LicenceID, input.TeamID)
	// There is an not found error here.
	if err != nil {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, err
	}

	// If this licence is not available to grant.
	if !licence.IsAvailable() {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, ErrLicenceUnavailable
	}

	// If another reader is already invited to accept this licence.
	// Admin should first revoke the invitation before invite another reader.
	if !licence.LastInviteeEmail.Valid && licence.LastInviteeEmail.String != input.Email {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, ErrInviteeMismatch
	}

	// Try to find the reader account by email.
	// Not found should not be considered an error here.
	invitee, err := env.FindReader(input.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, err
	}

	// If this reader has a valid membership, disallow
	// granting a new licence.
	if !invitee.Membership.IsExpired() {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, ErrAlreadyMember
	}

	if invitee.FtcID.IsZero() {
		invitee.Email = null.StringFrom(input.Email)
	}

	// Create Invitation instance based on the input data.
	inv, err := input.NewInvitation()
	if err != nil {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, err
	}

	// Update licence with by setting last_invitation column.
	baseLicence := licence.WithInvitation(inv)
	err = tx.SetLicenceInvited(baseLicence)
	if err != nil {
		_ = tx.Rollback()
		return admin.InvitedLicence{}, err
	}

	// Save the invitation
	err = tx.SaveInvitation(inv)
	if err != nil {
		return admin.InvitedLicence{}, err
	}

	if err := tx.Commit(); err != nil {
		return admin.InvitedLicence{}, err
	}

	return admin.InvitedLicence{
		Invitation: inv,
		Licence:    baseLicence,
		Plan:       licence.Plan,
		Assignee:   admin.Assignee{},
	}, nil
}

func (env Env) RevokeInvitation(invID, teamID string) error {
	tx, err := env.beginInvTx()
	if err != nil {
		return err
	}

	// Retrieve the invitation
	inv, err := tx.RetrieveInvitation(invID, teamID)
	// The invitation might not found.
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Retrieve the licence
	licence, err := tx.RetrieveLicence(inv.LicenceID, teamID)
	// Ignore the not found error
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return err
	}

	// Revoke the invitation
	inv = inv.Revoke()
	err = tx.RevokeInvitation(inv)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if licence.IsInvitationRevokable(inv.ID) {
		err := tx.UnlinkLicenceInvitation(licence)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// List invitations shows a list of invitations for a team.
func (env Env) ListInvitations(teamID string, page gorest.Pagination) ([]admin.Invitation, error) {
	var invs = make([]admin.Invitation, 0)

	err := env.db.Select(&invs, stmt.ListInvitation, teamID, page.Limit, page.Offset())

	if err != nil {
		return nil, err
	}

	return invs, nil
}

type InvitationsAndError struct {
	Invitations []admin.Invitation
	Error       error
}

func (env Env) AsyncListInvitations(teamID string, page gorest.Pagination) <-chan InvitationsAndError {
	r := make(chan InvitationsAndError)

	go func() {
		defer close(r)

		inv, err := env.ListInvitations(teamID, page)

		r <- InvitationsAndError{
			Invitations: inv,
			Error:       err,
		}
	}()

	return r
}

func (env Env) CountInvitation(teamID string) (int64, error) {
	var total int64

	err := env.db.Select(&total, stmt.CountInvitation, teamID)

	if err != nil {
		return total, err
	}

	return total, nil
}

type CountAndError struct {
	Total int64
	Error error
}

func (env Env) AsyncCountInvitation(teamID string) <-chan CountAndError {
	r := make(chan CountAndError)

	go func() {
		defer close(r)
		total, err := env.CountInvitation(teamID)

		r <- CountAndError{
			Total: total,
			Error: err,
		}
	}()

	return r
}
