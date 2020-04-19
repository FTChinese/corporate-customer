// Package stmt provides shared, reusable SQL statements.
// Most of them cannot be executed alone.
package stmt

const SignUp = `
INSERT INTO b2b.admin
SET id = :admin_id,
	email = :email,
	password_sha2 = UNHEX(SHA2(:password, 256)),
	vrf_token = UNHEX(:token),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const accountBase = `
SELECT a.id AS admin_id,
	a.team_id AS team_id,
	a.email AS email,
	a.display_name AS display_name,
	a.is_active AS active,
	a.verified AS verified`

// selectPassport retrieves admin's account
// and the team linked to it.
const selectPassport = accountBase + `,
	t.id AS team_id,
	t.name AS team_name
FROM b2b.admin AS a
	LEFT JOIN b2b.team AS t
	ON a.id = t.admin_id`

// PassportByAdminID retrieves Account + Team
// by account id.
const PassportByAdminID = selectPassport + `
WHERE a.id = ?
LIMIT 1`

const PassportByTeamID = selectPassport + `
WHERE t.id = ?
LIMIT 1`

// Login retrieves Account + Team by comparing
// credentials.
const Login = selectPassport + `
WHERE (a.email, a.password_sha2) = (?, UNHEX(SHA2(?, 256)))
LIMIT 1`

// selectAccount manipulate admin's account itself, without
// team data.
const selectAccount = accountBase + `
FROM b2b.admin AS a`

const AccountByID = selectAccount + `,
WHERE a.id = ?
LIMIT 1`

// AccountByVerifier retrieves admin account by
// verification token.
const AccountByVerifier = selectAccount + `
WHERE a.vrf_token = UNHEX(?)
LIMIT 1`

// EmailVerified set the verified column to true
// for a verification token.
const EmailVerified = `
UPDATE b2b.admin
	SET verified = 1
WHERE id = :admin_id
LIMIT 1`

const ReGenerateVrfToken = `
UPDATE b2b.admin
SET token = UNHEX(:token),
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT1`

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

const UpdateName = `
UPDATE b2b.admin
SET display_name = :display_name,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT 1`

const UpdatePassword = `
UPDATE b2b.admin
SET password_sha2 = UNHEX(SHA2(:password, 256))
	updated_utc = UTC_TIMESTAMP()
WHERE id = :amin_id
LIMIT 1`
