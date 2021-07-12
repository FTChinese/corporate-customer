package admin

const StmtInsertPwResetSession = `
INSERT INTO b2b.password_reset
SET token = UNHEX(:token),
	email = :email,
	source_url = :source_url,
	expires_in = :expires_in,
	created_utc = :created_utc`

// Do not removed the time comparison condition.
// It could reduce the chance of collision for app_code.
const selectPwResetSession = `

`

// StmtPwResetSessionByToken retrieves a password reset session
// by token for web app.
const StmtPwResetSessionByToken = `
SELECT LOWER(HEX(token)) AS token,
	email, 
	source_url,
	is_used,
	expires_in,
	created_utc
FROM b2b.password_reset
WHERE is_used = 0
	AND DATE_ADD(created_utc, INTERVAL expires_in SECOND) > UTC_TIMESTAMP()
	AND token = UNHEX(?)
LIMIT 1`

// StmtDisablePwResetToken flags a password reset token as invalid.
const StmtDisablePwResetToken = `
UPDATE b2b.password_reset
	SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`
