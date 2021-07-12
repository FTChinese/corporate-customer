package subs

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	admin2 "github.com/FTChinese/ftacademy/internal/pkg/admin"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
	"github.com/guregu/null"
)

// FindInvitationByToken tries to find an Invitation by token.
func (env Env) FindInvitationByToken(token string) (model2.Invitation, error) {
	var inv model2.Invitation
	err := env.dbs.Read.Get(&inv, stmt.InvitationByToken, token)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

// FindInvitedLicence tries to find a licence belong to
// an invitation.
func (env Env) FindInvitedLicence(claims model2.InviteeClaims) (model2.Licence, error) {
	var ls model2.LicenceSchema
	err := env.dbs.Read.Get(&ls, stmt.InvitedLicence, claims.LicenceID, claims.InvitationID)
	if err != nil {
		return model2.Licence{}, err
	}

	return ls.Licence()
}

// FindReader by email.
// Deprecated. Use API.
func (env Env) FindReader(email string) (model2.Reader, error) {
	var r model2.Reader
	err := env.dbs.Read.Get(&r, stmt.SelectReader, email)
	if err != nil {
		return r, err
	}

	r.Normalize()

	return r, nil
}

// CreateReader creates new FTC user.
// Deprecated. Use API.
func (env Env) CreateReader(s model2.SignUp) error {
	tx, err := env.dbs.Write.Beginx()
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
func (env Env) TakeSnapshot(snp model2.MemberSnapshot) error {
	_, err := env.dbs.Write.NamedExec(stmt.TakeSnapshot, snp)

	if err != nil {
		return err
	}

	return nil
}

// GrantLicence grants a licence to a reader.
func (env Env) GrantLicence(claims model2.InviteeClaims) (model2.InvitedLicence, error) {
	tx, err := env.beginInvTx()
	if err != nil {
		return model2.InvitedLicence{}, err
	}
	inv, err := tx.RetrieveInvitation(claims.InvitationID, claims.TeamID)
	// Not found
	if err != nil {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, err
	}
	if !inv.IsValid() {
		return model2.InvitedLicence{}, ErrInvalidInvitation
	}

	licence, err := tx.FindInvitedLicence(inv)
	// Not found
	if err != nil {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, err
	}
	// If licence cannot be granted, returns forbidden message.
	if !licence.CanBeGranted() {
		return model2.InvitedLicence{}, ErrLicenceTaken
	}

	mmb, err := tx.RetrieveMembership(claims.FtcID)
	if err != nil {
		_ = tx.Rollback()
		return model2.InvitedLicence{}, err
	}

	inv = inv.Accept()
	baseLicence := licence.GrantTo(claims.FtcID)
	newMmb := mmb.WithLicenceGranted(licence)

	// Create new membership based on licence
	if mmb.HasMembership() {
		err := tx.InsertMembership(newMmb)
		if err != nil {
			_ = tx.Rollback()
			return model2.InvitedLicence{}, err
		}
	} else {
		// Update current membership based on
		// licence.
		err := tx.UpdateMembership(newMmb)
		if err != nil {
			_ = tx.Rollback()
			return model2.InvitedLicence{}, err
		}

		// Back up.
		go func() {
			_ = env.TakeSnapshot(model2.NewMemberSnapshot(mmb))
		}()
	}

	if err := tx.LicenceGranted(baseLicence); err != nil {
		return model2.InvitedLicence{}, err
	}

	if err := tx.InvitationAccepted(inv); err != nil {
		return model2.InvitedLicence{}, err
	}

	if err := tx.Commit(); err != nil {
		return model2.InvitedLicence{}, err
	}

	// The returned data is used to compose a letter
	return model2.InvitedLicence{
		Invitation: inv,
		Licence:    baseLicence,
		Plan:       licence.Plan,
		Assignee: model2.Assignee{
			FtcID: null.StringFrom(claims.FtcID),
			Email: null.StringFrom(claims.Email),
		},
	}, nil
}

// FindInviteeOrg retrieves admin's data by team id.
// This is used to send admin an email after reader accepted
// an invitation
func (env Env) FindInviteeOrg(claims model2.InviteeClaims) (admin2.Passport, error) {
	var p admin2.Passport
	if err := env.dbs.Read.Get(&p, admin2.PassportByTeamID, claims.TeamID); err != nil {
		return admin2.Passport{}, err
	}

	return p, nil
}
