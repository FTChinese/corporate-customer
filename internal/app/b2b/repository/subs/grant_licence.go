package subs

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	"github.com/guregu/null"
)

// FindInvitationByToken tries to find an Invitation by token.
func (env Env) FindInvitationByToken(token string) (model.Invitation, error) {
	var inv model.Invitation
	err := env.db.Get(&inv, stmt.InvitationByToken, token)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

// FindInvitedLicence tries to find a licence belong to
// an invitation.
func (env Env) FindInvitedLicence(claims model.InviteeClaims) (model.Licence, error) {
	var ls model.LicenceSchema
	err := env.db.Get(&ls, stmt.InvitedLicence, claims.LicenceID, claims.InvitationID)
	if err != nil {
		return model.Licence{}, err
	}

	return ls.Licence()
}

func (env Env) FindReader(email string) (model.Reader, error) {
	var r model.Reader
	err := env.db.Get(&r, stmt.SelectReader, email)
	if err != nil {
		return r, err
	}

	r.Normalize()

	return r, nil
}

func (env Env) CreateReader(s model.SignUp) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(stmt.CreateReader, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(stmt.CreateReaderProfile, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(stmt.SaveReaderVrf, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// TakeSnapshot backs up a membership before
// modifying it.
func (env Env) TakeSnapshot(snp model.MemberSnapshot) error {
	_, err := env.db.NamedExec(stmt.TakeSnapshot, snp)

	if err != nil {
		return err
	}

	return nil
}

// GrantLicence grants a licence to a reader.
func (env Env) GrantLicence(claims model.InviteeClaims) (model.InvitedLicence, error) {
	tx, err := env.beginInvTx()
	if err != nil {
		return model.InvitedLicence{}, err
	}
	inv, err := tx.RetrieveInvitation(claims.InvitationID, claims.TeamID)
	// Not found
	if err != nil {
		_ = tx.Rollback()
		return model.InvitedLicence{}, err
	}
	if !inv.IsValid() {
		return model.InvitedLicence{}, ErrInvalidInvitation
	}

	licence, err := tx.FindInvitedLicence(inv)
	// Not found
	if err != nil {
		_ = tx.Rollback()
		return model.InvitedLicence{}, err
	}
	// If licence cannot be granted, returns forbidden message.
	if !licence.CanBeGranted() {
		return model.InvitedLicence{}, ErrLicenceTaken
	}

	mmb, err := tx.RetrieveMembership(claims.FtcID)
	if err != nil {
		_ = tx.Rollback()
		return model.InvitedLicence{}, err
	}

	inv = inv.Accept()
	baseLicence := licence.GrantTo(claims.FtcID)
	newMmb := mmb.WithLicenceGranted(licence)

	// Create new membership based on licence
	if mmb.HasMembership() {
		err := tx.InsertMembership(newMmb)
		if err != nil {
			_ = tx.Rollback()
			return model.InvitedLicence{}, err
		}
	} else {
		// Update current membership based on
		// licence.
		err := tx.UpdateMembership(newMmb)
		if err != nil {
			_ = tx.Rollback()
			return model.InvitedLicence{}, err
		}

		// Back up.
		go func() {
			_ = env.TakeSnapshot(model.NewMemberSnapshot(mmb))
		}()
	}

	if err := tx.LicenceGranted(baseLicence); err != nil {
		return model.InvitedLicence{}, err
	}

	if err := tx.InvitationAccepted(inv); err != nil {
		return model.InvitedLicence{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.InvitedLicence{}, err
	}

	// The returned data is used to compose a letter
	return model.InvitedLicence{
		Invitation: inv,
		Licence:    baseLicence,
		Plan:       licence.Plan,
		Assignee: model.Assignee{
			FtcID: null.StringFrom(claims.FtcID),
			Email: null.StringFrom(claims.Email),
		},
	}, nil
}

// FindInviteeOrg retrieves admin's data by team id.
// This is used to send admin an email after reader accepted
// an invitation
func (env Env) FindInviteeOrg(claims model.InviteeClaims) (model.Passport, error) {
	var p model.Passport
	if err := env.db.Get(&p, stmt.PassportByTeamID, claims.TeamID); err != nil {
		return model.Passport{}, err
	}

	return p, nil
}
