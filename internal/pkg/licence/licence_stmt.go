package licence

const colLicence = `
SELECT l.id AS licence_id,
	l.tier AS tier,
	l.cycle AS cycle,
	l.creator_id AS creator_id,
	l.team_id AS team_id,
	l.current_status AS lic_status,
	l.current_period_start_utc AS current_period_start_utc,
	l.current_period_end_utc AS current_period_end_utc,
	l.start_date_utc AS start_date_utc,
	l.trial_start_utc AS trial_start_utc,
	l.trial_end_utc AS trial_end_utc,
	l.latest_order_id AS latest_order_id,
	l.latest_price AS latest_price,
	l.latest_invitation AS latest_invitation,
	l.assignee_id AS assignee_id,
	l.created_utc AS created_utc,
	l.updated_utc AS updated_utc
`

const selectLicence = colLicence + `,
` + colAssigneeJSON + `
FROM b2b.licence AS l
	LEFT JOIN cmstmp01.userinfo AS a
	ON l.assignee_id = a.user_id
`

const StmtLicence = selectLicence + `
WHERE l.id = ? AND l.team_id = ?
LIMIT 1`

const StmtListLicences = selectLicence + `
WHERE l.team_id = ?
ORDER BY l.created_utc DESC
LIMIT ? OFFSET ?`

// StmtCountLicence is used to support pagination.
const StmtCountLicence = `
SELECT COUNT(*) AS total_licence
FROM b2b.licence
WHERE team_id = ?`

// StmtLockLicence retrieve a row from licence table and lock it for update.
const StmtLockLicence = colLicence + `
FROM b2b.licence AS l
WHERE l.id = ? AND l.team_id = ?
LIMIT 1
FOR UPDATE`

const StmtUpdateLicenceStatus = `
UPDATE b2b.licence
SET current_status = :lic_status,
	latest_invitation = :latest_invitation,
	assignee_id = :assignee_id,
	updated_utc = :updated_utc
WHERE id = :licence_id 
	AND team_id = :team_id
LIMIT 1`
