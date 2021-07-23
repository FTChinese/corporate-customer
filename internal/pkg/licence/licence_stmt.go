package licence

const colLicence = `
l.id AS licence_id,
l.tier AS tier,
l.cycle AS cycle,
l.creator_id AS creatorId,
l.team_id AS team_id,
l.current_period_start_utc AS current_period_start_utc,
l.current_period_end_utc AS current_period_end_utc,
l.start_date_utc AS start_date_utc,
l.trial_start_utc AS trial_start_utc,
l.trial_end_utc AS trial_end_utc,
l.latest_order_id AS latest_order_id,
l.latest_price AS latest_price,
l.current_status AS lic_status,
l.latest_invitation AS latest_invitation,
l.assignee_id AS assignee_id,
l.created_utc AS created_utc,
l.updated_utc AS updated_utc,
`

const selectLicence = `
SELECT ` + colLicence + `,
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
SET current_status = :current_status,
	latest_invitation = :latest_invitation
	assignee_id = :assignee_id
	updated_utc = :updated_utc
WHERE id = :licence_id 
	AND team_id = :team_id
LIMIT 1`

// StmtLicenceInvited updates the current_status
// and last_invitation columns after an invitation is sent
// Deprecated
const StmtLicenceInvited = `
UPDATE b2b.licence
SET current_status = :current_status,
	latest_invitation = :latest_invitation
	updated_utc = :updated_utc
WHERE id = :licence_id 
	AND team_id = :team_id
LIMIT 1`

// StmtInvitedLicence retrieves a licence for an invitation.
// Deprecated
const StmtInvitedLicence = selectLicence + `
WHERE l.latest_invitation_id = ?
LIMIT 1`

// StmtLicenceByID retrieves a single row.
// Deprecated.
const StmtLicenceByID = selectLicence + `
WHERE l.id = ? AND team_id = ?
LIMIT 1`

// LockLicenceByID locks and retrieves a row of licence
// when creating an invitation for it.
// Deprecated.
const LockLicenceByID = StmtLicenceByID + `
FOR UPDATE`

// StmtLockInvitedLicence locks a licence row belong to an
// invitation.
// Deprecated.
const StmtLockInvitedLicence = StmtInvitedLicence + `
FOR UPDATE`

// SetLicenceGranted after user accepted invitation.
// Deprecated.
const SetLicenceGranted = `
UPDATE b2b.licence
SET assignee_id = :assignee_id,
	current_status = :current_status,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :licence_id
	AND team_id = :team_id
LIMIT 1`

// SetLicenceRevoked if admin decides to revoke the the licence
// from a user.
// User's membership status should also be synced.
// Deprecated.
const SetLicenceRevoked = `
UPDATE b2b.licence
SET assignee_id = NULL,
	current_status = :current_status,
	last_invitation_id = NULL,
	last_invitee_email = NULL,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :licence_id
	AND team_id = :team_id
LIMIT 1`
