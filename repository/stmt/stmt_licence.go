package stmt

const licenceCols = `l.id AS licence_id,
l.team_id,
l.plan_id,
l.expire_date,
l.assignee_id
l.is_active,
l.created_utc AS created_utc,
l.updated_utc AS updated_utc`

const selectExpandedLicence = `
SELECT ` + licenceCols + `,
` + planCols + `,
` + readerAccountCols +
	`FROM b2b.licence AS l
LEFT JOIN subs.plan AS p
	ON l.plan_id = p.id
LEFT JOIN cmstmp01.userinfo AS u
	ON l.assignee_id = u.user_id`

// Select a single licence belonging to a team.
const ExpandedLicence = selectExpandedLicence + `
WHERE l.id = ? AND l.team_id = ?
LIMIT 1`

// FindExpandedLicence searches a licence by id.
// Used when verifying a user's invitation.
const FindExpandedLicence = selectExpandedLicence + `
WHERE l.id = ?
LIMIT 1`

// Select a list of licence for a team.
const ListExpandedLicences = selectExpandedLicence + `
WHERE l.team_id = ?
ORDER BY l.created_utc DESC
LIMIT ? OFFSET ?`

const CountLicence = `
SELECT COUNT(*) AS total_licence
FROM b2b.licence
WHERE team_id = ?`

// LockLicence locks a row of licence
// when granting it to user.
const LockLicence = `
SELECT ` + licenceCols + `
FROM b2b.licence AS l
WHERE l.id = ?
LIMIT 1
FOR UPDATE`
