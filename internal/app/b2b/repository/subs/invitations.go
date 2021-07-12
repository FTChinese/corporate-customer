package subs

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
	"github.com/FTChinese/go-rest"
	"github.com/guregu/null"
)

// CreateInvitation creates a new invitation for a licence.
// To create an invitation letter, we need the following
// information:
// * Assignee
// * Invitation.Token
// * Plan
func (env Env) CreateInvitation(inv model2.Invitation) (model2.InvitedLicence, error) {
	tx, err := env.beginInvTx()
	if err != nil {
		return model2.InvitedLicence{}, err
	}

	// Retrieve the licence.
	licence, err := tx.RetrieveLicence(inv.LicenceID, inv.TeamID)
	// There is an not found error here.
	if err != nil {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, err
	}

	// If this licence is not available to grant.
	if !licence.IsAvailable() {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, ErrLicenceUnavailable
	}

	// If another reader is already invited to accept this licence.
	// Admin should first revoke the invitation before invite another reader.
	if !licence.LastInviteeEmail.Valid && licence.LastInviteeEmail.String != inv.Email {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, ErrInviteeMismatch
	}

	// Try to find the reader account by email.
	// Not found should not be considered an error here.
	invitee, err := env.FindReader(inv.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, err
	}

	// If this reader has a valid membership, disallow
	// granting a new licence.
	if !invitee.Membership.IsExpired() {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, ErrAlreadyMember
	}

	if invitee.FtcID.IsZero() {
		invitee.Email = null.StringFrom(inv.Email)
	}

	// Update licence with by setting last_invitation column.
	baseLicence := licence.Invited(inv)
	err = tx.SetLicenceInvited(baseLicence)
	if err != nil {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, err
	}

	// Save the invitation
	err = tx.SaveInvitation(inv)
	if err != nil {
		return model2.InvitedLicence{}, err
	}

	if err := tx.Commit(); err != nil {
		return model2.InvitedLicence{}, err
	}

	return model2.InvitedLicence{
		Invitation: inv,
		Licence:    baseLicence,
		Plan:       licence.Plan,
		Assignee:   invitee.Assignee,
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
	licence, err := tx.FindInvitedLicence(inv)
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

	if licence.CanInvitationBeRevoked() {
		err := tx.UnlinkInvitedLicence(licence)
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
func (env Env) ListInvitations(teamID string, page gorest.Pagination) ([]model2.Invitation, error) {
	var invs = make([]model2.Invitation, 0)

	err := env.dbs.Read.Select(&invs, stmt.ListInvitation, teamID, page.Limit, page.Offset())

	if err != nil {
		return nil, err
	}

	return invs, nil
}

func (env Env) AsyncListInvitations(teamID string, page gorest.Pagination) <-chan model2.InvitationList {
	r := make(chan model2.InvitationList)

	go func() {
		defer close(r)

		inv, err := env.ListInvitations(teamID, page)

		r <- model2.InvitationList{
			Data: inv,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountInvitation(teamID string) (int64, error) {
	var total int64

	err := env.dbs.Read.Select(&total, stmt.CountInvitation, teamID)

	if err != nil {
		return total, err
	}

	return total, nil
}

func (env Env) AsyncCountInvitation(teamID string) <-chan model2.InvitationList {
	r := make(chan model2.InvitationList)

	go func() {
		defer close(r)
		total, err := env.CountInvitation(teamID)

		r <- model2.InvitationList{
			Total: total,
			Err:   err,
		}
	}()

	return r
}
