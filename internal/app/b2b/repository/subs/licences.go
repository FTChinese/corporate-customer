package subs

import (
	"github.com/FTChinese/b2b/internal/app/b2b/model"
	"github.com/FTChinese/b2b/internal/app/b2b/stmt"
	gorest "github.com/FTChinese/go-rest"
)

// RevokeLicence revokes a licence granted to a reader.
// If Licence.Status is not LicStatusGranted, no-ops will
// be performed.
// For a licence waiting for its invitation accepted,
// use RevokeInvitation instead of this one.
func (env Env) RevokeLicence(id, teamID string) error {
	tx, err := env.beginInvTx()
	if err != nil {
		return err
	}

	licence, err := tx.RetrieveLicence(id, teamID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// If the licence is not granted yet.
	if licence.Status != model.LicStatusGranted {
		_ = tx.Rollback()
		return nil
	}

	// Get assignee current membership
	mmb, err := tx.RetrieveMembership(licence.AssigneeID.String)
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

	// Update licence data.
	baseLicence := licence.InvitationRevoked()
	err = tx.RevokeLicence(baseLicence)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// LoadExpLicence retrieves a licence, together with its
// subscription plan and the user to whom it was assigned.
// If the licence is not assigned yet, assignee fields are empty..
func (env Env) LoadExpLicence(id, teamID string) (model.ExpandedLicence, error) {
	var ls model.ExpLicenceSchema
	err := env.db.Get(&ls, stmt.ExpandedLicence, id, teamID)

	if err != nil {
		return model.ExpandedLicence{}, err
	}

	return ls.ExpandedLicence()
}

// ListExpLicence shows a list all licence.
// Each licence's plan, invitation, assignee are attached.
func (env Env) ListExpLicence(teamID string, page gorest.Pagination) ([]model.ExpandedLicence, error) {
	var ls = make([]model.ExpLicenceSchema, 0)

	err := env.db.Select(&ls, stmt.ListExpandedLicences, teamID, page.Limit, page.Offset())

	if err != nil {
		return nil, err
	}

	el := make([]model.ExpandedLicence, 0)
	for _, v := range ls {
		l, err := v.ExpandedLicence()
		if err != nil {
			return nil, err
		}
		el = append(el, l)
	}
	return el, nil
}

func (env Env) AsyncListExpLicence(teamID string, page gorest.Pagination) <-chan model.PagedExpLicences {
	r := make(chan model.PagedExpLicences)

	go func() {
		defer close(r)
		licences, err := env.ListExpLicence(teamID, page)

		r <- model.PagedExpLicences{
			Data: licences,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountLicences(teamID string) (int64, error) {
	var total int64
	if err := env.db.Get(&total, stmt.CountLicence, teamID); err != nil {
		return total, err
	}

	return total, nil
}

func (env Env) AsyncCountLicences(teamID string) <-chan model.PagedExpLicences {
	r := make(chan model.PagedExpLicences)

	go func() {
		defer close(r)
		total, err := env.CountLicences(teamID)

		r <- model.PagedExpLicences{
			Total: total,
			Err:   err,
		}
	}()

	return r
}
