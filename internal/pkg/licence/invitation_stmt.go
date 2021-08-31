package licence

const StmtCreateInvitation = `
INSERT INTO b2b.invitation
SET id = :invite_id,
	admin_id = :admin_id,
	team_id = :team_id,
	description = :invite_desc,
	expiration_days = :invite_expiration_days,
	invitee_email = :invite_email,
	licence_id = :licence_id,
	current_status = :invite_status,
	token = UNHEX(:invite_token),
	created_utc = :created_utc,
	updated_utc = :updated_utc`

// colInvitation does not include the token column
// as it is only visible to the user who received invitation
// email.
const colInvitation = `
SELECT i.id AS invite_id,
	i.admin_id AS admin_id,
	i.team_id AS team_id,
	i.current_status AS invite_status,
	i.description AS invite_desc,
	i.expiration_days AS invite_expiration_days,
	i.invitee_email AS invite_email,
	i.licence_id AS licence_id,
	LOWER(HEX(i.token)) AS invite_token,
	i.created_utc AS created_utc,
	i.updated_utc AS updated_utc
FROM b2b.invitation AS i
`

const StmtInvitationByID = colInvitation + `
WHERE i.id = ? AND i.team_id = ?
LIMIT 1
`

// StmtLockInvitation locks and retrieves a row of invitation
// upon revoking.
const StmtLockInvitation = StmtInvitationByID + `
FOR UPDATE`

// StmtInvitationByToken retrieves an expanded invitation
// by the token sent in an email.
// This is used to verify an invitation email.
const StmtInvitationByToken = colInvitation + `
WHERE i.token = UNHEX(?)
LIMIT 1`

// StmtListInvitation shows a list
// of invitations with its assignee attached.
// This is used by admin to views invitations it sent.
const StmtListInvitation = colInvitation + `
WHERE i.team_id = ?
ORDER BY i.created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountInvitation = `
SELECT COUNT(*)
FROM b2b.invitation
WHERE team_id = ?`

const StmtUpdateInvitationStatus = `
UPDATE b2b.invitation
SET current_status = :invite_status,
	updated_utc = :updated_utc
WHERE id = :invite_id
	AND team_id = :team_id
LIMIT 1`
