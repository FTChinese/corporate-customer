package licence

const colLicence = `
l.id AS licence_id,
l.team_id AS team_id,
l.created_utc AS created_utc,
l.current_status AS licence_status,
l.current_period_start_utc AS current_period_start_utc,
l.current_period_end_utc AS current_period_end_utc,
l.start_date_utc AS start_date_utc,
l.trial_start_utc AS trial_start_utc,
l.trial_end_utc AS trial_end_utc,
l.latest_invoice_id AS latest_invoice_id,
l.updated_utc AS updated_utc,
l.latest_price AS latest_price,
l.latest_invitation AS latest_invitation
`

// Retrieve Assignee as JSON object so that we don't need to
// create extra types to convert string to JSON.
const colAssignee = `
JSON_OBJECT(
	"ftcId", a.user_id,
	"email", a.email,
	"userName", a.user_name
) AS assignee`

const selectLicence = `
SELECT ` + colLicence + `,
` + colAssignee + `
FROM b2b.licence AS l
	LEFT JOIN cmstmp01.userinfo AS a
	ON l.assignee_id = a.user_id
`

const StmtLicence = selectLicence + `
WHERE l.id = ? AND l.team_id = ?
LIMIT 1`

// StmtInvitedLicence retrieves a licence for an invitation.
const StmtInvitedLicence = selectLicence + `
WHERE l.latest_invitation_id = ?
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

// LicenceByID retrieves a single row.
const LicenceByID = selectLicence + `
WHERE l.id = ? AND team_id = ?
LIMIT 1`

// LockLicenceByID locks and retrieves a row of licence
// when creating an invitation for it.
const LockLicenceByID = LicenceByID + `
FOR UPDATE`

// InvitedLicence locks a licence row belong to an
// invitation.
const LockInvitedLicence = StmtInvitedLicence + `
FOR UPDATE`

// SetLicenceInvited updates the current_status
// and last_invitation columns after an invitation is sent
const SetLicenceInvited = `
UPDATE b2b.licence
SET current_status = :current_status,
	last_invitation_id = :last_invitation_id,
	last_invitee_email = :last_invitee_email,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :licence_id 
	AND team_id = :team_id
LIMIT 1`

// RevokeLicenceInvitation removes last_invitation when
// admin revokes an invitation.
const RevokeLicenceInvitation = `
UPDATE b2b.licence
SET current_status = :current_status,
	last_invitation_id = NULL,
	last_invitee_email = NULL,
	updated_utc = UTC_TIMESTAMP
WHERE id = :licence_id 
	AND team_id = :team_id
LIMIT 1`

// SetLicenceGranted after user accepted invitation.
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
