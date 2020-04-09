package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

const stmtCreateInvitation = `
INSERT INTO b2b.invitation
SET id = :invitation_id,
	licence_id = :licence_id,
	token = :UNHEX(:token),
	invitee_email = :invitee_email,
	description = :description,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

func (env Env) CreateInvitation(inv admin.Invitation) error {
	_, err := env.db.NamedExec(stmtCreateInvitation, inv)

	if err != nil {
		return err
	}

	return nil
}

// List invitations shows a list of invitations for a team.
func (env Env) ListInvitations(teamID string) ([]admin.ExpandedInvitation, error) {
	var invs = make([]admin.InvitationSchema, 0)

	err := env.db.Select(&invs, stmt.ListExpandedInvitation, teamID)

	if err != nil {
		return nil, err
	}

	eis := make([]admin.ExpandedInvitation, 0)
	for _, v := range invs {
		eis = append(eis, v.ExpandedInvitation())
	}

	return eis, nil
}

const stmtRevokeInvitation = `
UPDATE b2b.invitation
SET revoked = 1
WHERE id = ?
	AND accepted = 0
	AND team_id = ?
LIMIT 1`

func (env Env) RevokeInvitation(id, teamID string) error {
	_, err := env.db.Exec(stmtRevokeInvitation, id, teamID)

	if err != nil {
		return err
	}

	return nil
}

// LoadInvitation shows a single invitation to admin.
func (env Env) LoadInvitation(id, teamID string) (admin.ExpandedInvitation, error) {
	var i admin.InvitationSchema
	err := env.db.Get(&i, stmt.ExpandedInvitation, id, teamID)

	if err != nil {
		return admin.ExpandedInvitation{}, err
	}

	return i.ExpandedInvitation(), nil
}
