package licence

// DO not change the order the team_id in the WHERE clause
// because the defined the PRIMARY KEY as (team_id, email)
// for this table.
// See https://dev.mysql.com/doc/refman/8.0/en/multiple-column-indexes.html.

const InsertStaffer = `
INSERT IGNORE INTO b2b.staff
SET email = :email,
	ftc_id = :ftc_id,
	team_id = :team_id,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const DeleteStaffer = `
DELETE FROM b2b.staff
WHERE team_id = :team_id
	AND email = :email
LIMIT 1`

const SetStaffFtcID = `
UPDATE b2b.staff
SET ftc_id = :ftc_id
WHERE team_id = :team_id
	AND email = :email
	AND ftc_id = NULL
LIMIT 1`

const ListStaff = `
SELECT id,
	email,
	ftc_id,
	team_id,
	created_utc,
	updated_utc
FROM b2b.staff
WHERE team_id = ?
ORDER BY email ASC
LIMIT ? OFFSET ?`

const CountStaff = `
SELECT COUNT(*)
FROM b2b.staff
WHERE team_id = ?`
