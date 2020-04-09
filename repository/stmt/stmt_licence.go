package stmt

import "github.com/FTChinese/b2b/models/builder"

const licenceCols = `l.id AS licence_id,
l.team_id,
l.plan_id,
l.expire_date,
l.assignee_id
l.is_active,
l.created_utc AS created_utc,
l.updated_utc AS updated_utc`

var selectExpandedLicence = builder.NewSelect().
	AddRawColumn(licenceCols).
	AddRawColumn(planCols).
	AddRawColumn(readerAccountCols).
	From(`b2b.licence AS l
LEFT JOIN subs.plan AS p
	ON l.plan_id = p.id
LEFT JOIN cmstmp01.userinfo AS u
	ON l.assignee_id = u.user_id`)

// Select a single licence belonging to a team.
var ExpandedLicence = selectExpandedLicence.
	Where(`l.id = ? AND l.team_id = ?`).
	Limit(1).
	Build()

// FindExpandedLicence searches a licence by id.
// Used when verifying a user's invitation.
var FindExpandedLicence = selectExpandedLicence.
	Where("l.id = ?").
	Limit(1).
	Build()

// Select a list of licence for a team.
var ListExpandedLicences = selectExpandedLicence.
	Where("l.team_id = ?").
	OrderBy("l.created_utc DESC").
	Paged().
	Build()

// LockLicence locks a row of licence
// when granting it to user.
var LockLicence = builder.NewSelect().
	AddRawColumn(licenceCols).
	From("b2b.licence AS l").
	Where("l.id = ?").
	Limit(1).
	Lock().
	Build()
