package admin

const StmtInsertPwResetSession = `
INSERT INTO b2b.password_reset
SET token = UNHEX(:token),
	email = :email,
	expires_in = :expires_in,
	created_utc = :created_utc`

// StmtPwResetSessionByToken retrieves a password reset session
// by token for web app.
const StmtPwResetSessionByToken = `
SELECT LOWER(HEX(token)) AS token,
	email,
	is_used,
	expires_in,
	created_utc
FROM b2b.password_reset
WHERE is_used = 0
	AND token = UNHEX(?)
LIMIT 1`

// StmtDisablePwResetToken flags a password reset token as invalid.
const StmtDisablePwResetToken = `
UPDATE b2b.password_reset
SET is_used = :is_used,
	updated_utc = :updated_utc
WHERE token = UNHEX(:token)
LIMIT 1`
