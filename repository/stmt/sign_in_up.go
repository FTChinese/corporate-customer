package stmt

// LogIn verifies email and password combination.
const LogIn = `
SELECT id AS admin_id
FROM b2b.admin
WHERE (email, password_sha2) = (?, UNHEX(SHA2(?, 256)))
LIMIT 1`

// SignUp inserts an admin account.
const SignUp = `
INSERT INTO b2b.admin
SET id = :admin_id,
	email = :email,
	password_sha2 = UNHEX(SHA2(:password, 256)),
	vrf_token = UNHEX(:vrf_token),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

// ReGeneratedVrfToken regenerate verification token upon request.
// An email should also be sent.
const ReGenerateVrfToken = `
UPDATE b2b.admin
SET token = UNHEX(:token),
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT1`

// SetEmailVerified set the verified column to true
// for a verification token.
const SetEmailVerified = `
UPDATE b2b.admin
	SET verified = TRUE
WHERE vrf_token = :vrf_token
LIMIT 1`
