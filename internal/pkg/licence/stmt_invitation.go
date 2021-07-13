package licence

// selectInvitation does not include the token column
// as it is only visible to the user who received invitation
// email.
const selectInvitation = `
SELECT i.id AS invitation_id,
	i.licence_id AS licence_id,
	i.team_id AS team_id,
	i.invitee_email AS email,
	i.description AS description,
	i.expiration_days AS expiration_days,
	i.current_status AS current_status,
	i.created_utc AS inv_created_utc,
	i.updated_utc AS inv_updated_utc
FROM b2b.invitation AS i`

// LockInvitationByID locks and retrieves a row of invitation
// upon revoking.
const LockInvitationByID = selectInvitation + `
WHERE i.id = ? AND team_id = ?
LIMIT 1
FOR UPDATE`

// InvitationByToken retrieves an expanded invitation
// by the token sent in an email.
// This is used to verify an invitation email.
const InvitationByToken = selectInvitation + `
WHERE i.token = UNHEX(?)
LIMIT 1`

// LockInvitationByToken searches for an invitation by token
// and lock it when granting licence.
const LockInvitationByToken = InvitationByToken + `
FOR UPDATE`

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

const CreateInvitation = `
INSERT INTO b2b.invitation
SET id = :invitation_id,
	licence_id = :licence_id,
	token = UNHEX(:token),
	invitee_email = :invitee_email,
	description = :description,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

// RevokeInvitation mark an invitation as invalid.
// If an invitation is already accepted, this operation
// does nothing.
const RevokeInvitation = `
UPDATE b2b.invitation
SET current_status = :current_status,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :invitation_id
	AND team_id = :team_id
LIMIT 1`

// AcceptInvitation for a reader who received invitation email.
const AcceptInvitation = `
UPDATE b2b.invitation
SET current_status = :current_status,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :invitation_id
LIMIT 1`
