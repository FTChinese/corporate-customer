package stmt

const InsertTeamMember = `
INSERT IGNORE INTO b2b.team_member
SET email = :email,
	ftc_id = :ftc_id,
	team_id = :team_id,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const DeleteTeamMember = `
DELETE FROM b2b.team_member
WHERE id = :id
	AND team_id = :team_id
LIMIT 1`

const SetTeamMemberFtcID = `
UPDATE b2b.team_member
SET ftc_id = :ftc_id
WHERE email = :email
	AND ftc_id = NULL
LIMIT 1`

const ListTeamMembers = `
SELECT id,
	email,
	ftc_id,
	team_id,
	created_utc,
	updated_utc
FROM b2b.team_member
WHERE team_id = ?
ORDER BY email ASC
LIMIT ? OFFSET ?`

const CountTeamMembers = `
SELECT COUNT(*)
FROM b2b.team_member
WHERE team_id = ?`
