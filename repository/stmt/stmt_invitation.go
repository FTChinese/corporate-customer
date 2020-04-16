package stmt

const CreateInvitation = `
INSERT INTO b2b.invitation
SET id = :invitation_id,
	licence_id = :licence_id,
	token = UNHEX(:token),
	invitee_email = :invitee_email,
	description = :description,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const RevokeInvitation = `
UPDATE b2b.invitation
SET revoked = 1
WHERE id = ?
	AND accepted = 0
	AND team_id = ?
LIMIT 1`

// selectInvitation does not include the token column
// as it is only visible to the user who received invitation
// email.
const selectInvitation = `
SELECT i.id AS invitation_id,
	i.licence_id AS licence_id,
	i.team_id AS team_id,
	i.invitee_email AS email,
	i.expiration_days AS expiration_days,
	i.description AS description,
	i.accepted AS accepted,
	i.revoked AS revoked,
	i.created_utc AS inv_created_utc,
	i.updated_utc AS inv_updated_utc
FROM b2b.invitation AS i`

// ListExpandedInvitation shows a list
// of invitations with its assignee attached.
// This is used by admin to views invitations it sent.
const ListInvitation = selectInvitation + `
WHERE i.team_id = ?
ORDER BY i.created_utc DESC
LIMIT ? OFFSET ?`

const CountInvitation = `
SELECT COUNT(*)
FROM b2b.invitation
WHERE team_id = ?`

// FindExpandedInvitation retrieves an expanded invitation
// by the token sent in an email.
// This is used to verify an invitation email.
const FindInvitation = selectInvitation + `
WHERE i.token = UNHEX(?)
LIMIT 1`

// LockInvitation locks a row in invitation
// when granting user a licence.
var LockInvitation = selectInvitation + `
WHERE i.id = ?
LIMIT 1
FOR UPDATE`
