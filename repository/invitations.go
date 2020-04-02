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

const stmtListInvitation = stmt.Invitation + `
WHERE team_id = ?
ORDER BY created_utc DESC`

func (env Env) ListInvitations(teamID string) ([]admin.Invitation, error) {
	var invs = make([]admin.Invitation, 0)

	err := env.db.Select(&invs, stmtListInvitation, teamID)

	if err != nil {
		return invs, err
	}

	return invs, nil
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

const stmtInvitation = `
WHERE id = ?
	AND team_id = ?
LIMIT 1`

// LoadInvitation shows a single invitation to admin.
func (env Env) LoadInvitation(id, teamID string) (admin.Invitation, error) {
	var i admin.Invitation
	err := env.db.Get(&i, stmtInvitation, id, teamID)

	if err != nil {
		return i, err
	}

	return i, nil
}
