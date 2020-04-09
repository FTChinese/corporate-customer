package stmt

import "github.com/FTChinese/b2b/models/builder"

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

var selectExpInvtBuilder = builder.NewSelect().
	AddColumn(invitationCols).
	AddColumn(readerAccountCols).
	From(`b2b.invitation AS i
LEFT JOIN cmstmp01.userinfo AS u
	ON i.invitee_email = u.email`)

// ListExpandedInvitation shows a list
// of invitations with its assignee attached.
// This is used by admin to views invitations it sent.
var ListExpandedInvitation = selectExpInvtBuilder.
	Where("i.team_id = ?").
	OrderBy("i.created_utc DESC").
	Paged().
	Build()

// ExpandedInvitation retrieves an expanded invitation
// belonging to a team
var ExpandedInvitation = selectExpInvtBuilder.
	Where("i.id = ? ADN i.team_id = ?").
	Limit(1).
	Build()

// FindExpandedInvitation retrieves an expanded invitation
// by the token sent in an email.
// This is used to verify an invitation email.
var FindExpandedInvitation = selectExpInvtBuilder.
	Where("i.token = UNHEX(?)").
	Limit(1).
	Build()

// LockInvitation locks a row in invitation
// when granting user a licence.
var LockInvitation = builder.NewSelect().
	AddColumn(invitationCols).
	From(`b2b.invitation AS i`).
	Where("i.id = ?").
	Limit(1).
	Lock().
	Build()
