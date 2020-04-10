package stmt

const invitationCols = `
i.id AS invitation_id,
i.licence_id AS licence_id,
i.team_id AS team_id,
LOWER(HEX(i.token)) AS token,
i.description AS description,
i.accepted AS accepted,
i.revoked AS revoked,
i.created_utc AS inv_created_utc,
i.updated_utc AS inv_updated_utc`

const selectExpInvtBuilder = `
SELECT ` + invitationCols + `,
` + readerAccountCols + `
FROM b2b.invitation AS i
LEFT JOIN cmstmp01.userinfo AS u
	ON i.invitee_email = u.email`

// ListExpandedInvitation shows a list
// of invitations with its assignee attached.
// This is used by admin to views invitations it sent.
const ListExpandedInvitation = selectExpInvtBuilder + `
WHERE i.team_id = ?
ORDER BY i.created_utc DESC
LIMIT ? OFFSET ?`

const CountInvitation = `
SELECT COUNT(*)
FROM b2b.invitation
WHERE team_id = ?`

// ExpandedInvitation retrieves an expanded invitation
// belonging to a team
const ExpandedInvitation = selectExpInvtBuilder + `
WHERE i.id = ? ADN i.team_id = ?
LIMIT 1`

// FindExpandedInvitation retrieves an expanded invitation
// by the token sent in an email.
// This is used to verify an invitation email.
const FindExpandedInvitation = selectExpInvtBuilder + `
WHERE i.token = UNHEX(?)
LIMIT 1`

// LockInvitation locks a row in invitation
// when granting user a licence.
var LockInvitation = `
SELECT ` + invitationCols + `
FROM b2b.invitation AS i
WHERE i.id = ?
LIMIT 1
FOR UPDATE`
