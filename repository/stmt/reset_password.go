package stmt

// InsertPasswordResetToken generates a new token
// to help resetting password.
const InsertPasswordResetToken = `
INSERT INTO b2b.password_reset
SET token = UNHEX(:token),
	email = :email,
	created_utc = UTC_TIMESTAMP()`

// RetrievePasswordReset retrieves an account for a
// password reset token.
const RetrievePasswordReset = retrieveAdmin + `,
	r.expires_in,
	r.created_utc
FROM b2b.password_reset AS r
	LEFT JOIN b2b.admin AS a
	ON r.email = a.email
WHERE r.token = UNHEX(?)
	AND r.is_used = FALSE
LIMIT 1`

// DeactivatePasswordResetToken set a token's is_used
// column to true so that it won't be retrieved in the future.
const DeactivatePasswordResetToken = `
UPDATE b2b.password_reset
	SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`
