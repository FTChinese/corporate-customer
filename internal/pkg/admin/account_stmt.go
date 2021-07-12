package admin

const StmtCreateAdmin = `
INSERT INTO b2b.admin
SET id = :admin_id,
	email = :email,
	password_sha2 = UNHEX(SHA2(:password, 256)),
	vrf_token = UNHEX(:token),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const colBaseAccount = `
SELECT a.id AS admin_id,
	t.id AS team_id,
	a.email AS email,
	a.display_name AS display_name,
	a.is_active AS active,
	a.verified AS verified
`

const selectBaseAccount = colBaseAccount + `
FROM b2b.admin AS a
	LEFT JOIN b2b.team AS t
	ON a.id = t.admin_id
`

const StmtBaseAccountByID = selectBaseAccount + `
WHERE a.id = ?
LIMIT 1`

const StmtBaseAccountByEmail = selectBaseAccount + `
WHERE a.email = ?
LIMIT 1`

// selectPassport retrieves admin's account
// and the team linked to it.
const selectPassport = colBaseAccount + `,
	t.id AS team_id,
	t.name AS team_name
FROM b2b.admin AS a
	LEFT JOIN b2b.team AS t
	ON a.id = t.admin_id`

// PassportByTeamID retrieve passport by team id.
// Deprecated.
const PassportByTeamID = selectPassport + `
WHERE t.id = ?
LIMIT 1`

const StmtUpdateName = `
UPDATE b2b.admin
SET display_name = :display_name,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT 1`

const StmtUpdatePassword = `
UPDATE b2b.admin
SET password_sha2 = UNHEX(SHA2(:password, 256))
	updated_utc = UTC_TIMESTAMP()
WHERE id = :amin_id
LIMIT 1`
