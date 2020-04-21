package stmt

const licenceCols = `l.id AS licence_id,
l.team_id AS team_id,
l.assignee_id AS assignee_id,
l.expire_date AS expire_date,
l.trial_start_date AS trial_start_date,
l.trial_end_date AS trial_end_date,
l.current_status AS current_status,
l.created_utc AS created_utc,
l.updated_utc AS updated_utc,
l.current_plan AS current_plan,
l.last_invitation_id AS last_invitation_id,
l.last_invitee_email AS last_invitee_email`

const selectLicence = `
SELECT ` + licenceCols + `
FROM b2b.licence AS l`

// LicenceByID retrieves a single row.
const LicenceByID = selectLicence + `
WHERE l.id = ? AND team_id = ?
LIMIT 1`

// LockLicenceByID locks and retrieves a row of licence
// when creating an invitation for it.
const LockLicenceByID = LicenceByID + `
FOR UPDATE`

// InvitedLicence retrieves a licence for an invitation.
const InvitedLicence = selectLicence + `
WHERE l.id = ? AND last_invitation_id = ?
LIMIT 1`

// InvitedLicence locks a licence row belong to an
// invitation.
const LockInvitedLicence = InvitedLicence + `
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

const selectExpandedLicence = `
SELECT ` + licenceCols + `,
` + readerAccountCols + `
FROM b2b.licence AS l
LEFT JOIN cmstmp01.userinfo AS u
	ON l.assignee_id = u.user_id`

// Select a single licence belonging to a team.
const ExpandedLicence = selectExpandedLicence + `
WHERE l.id = ? AND l.team_id = ?
LIMIT 1`

// Select a list of licence for a team.
const ListExpandedLicences = selectExpandedLicence + `
WHERE l.team_id = ?
ORDER BY l.created_utc DESC
LIMIT ? OFFSET ?`

// CountLicence is used to support pagination.
const CountLicence = `
SELECT COUNT(*) AS total_licence
FROM b2b.licence
WHERE team_id = ?`
