// Package stmt provides shared, reusable SQL statements.
// Most of them cannot be executed alone.
package stmt

const accountBase = `
SELECT a.id AS admin_id,
	a.team_id AS team_id,
	a.email AS email,
	a.display_name AS display_name,
	a.is_active AS active,
	a.verified AS verified`

const selectAccount = accountBase + `
FROM b2b.admin AS a`

const selectJWTAccount = accountBase + `,
	t.team_id AS team_id
FROM b2b.admin AS a
	LEFT JOIN b2b.team AS t
	ON a.id = t.admin_id`

const JWTAccount = selectJWTAccount + `
WHERE a.id = ?
LIMIT 1`

const Login = selectJWTAccount + `
WHERE (a.email, a.password_sha2) = (?, UNHEX(SHA2(?, 256)))
LIMIT 1`

const AccountByID = selectAccount + `,
WHERE a.id = ?
LIMIT 1`

const AccountByVerifier = selectAccount + `
WHERE a.vrf_token = UNHEX(?)
LIMIT 1`

// AccountByEmail retrieves an admin's account
// by email column.
// Used when requesting a password reset letter.
const AccountByEmail = selectAccount + `
WHERE email = ?
LIMIT`

const AccountForPwReset = accountBase + `
FROM b2b.password_reset AS r
	INNER JOIN b2b.admin AS a
	ON r.email = a.email
WHERE r.token = UNHEX(?)
	AND r.is_used = 0
	AND DATE_ADD(r.updated_utc, INTERVAL r.expires_in SECOND) > UTC_TIMESTAMP()
LIMIT 1`

const AdminProfile = accountBase + `,
a.created_utc AS created_utc,
a.updated_utc AS updated_utc
FROM b2b.admin AS a
WHERE a.id = ?
LIMIT 1`

const PasswordMatched = `
SELECT password_sha2 = UNHEX(SHA2(?, 256)) AS matched
FROM b2b.admin
WHERE id = ?
LIMIT 1`
