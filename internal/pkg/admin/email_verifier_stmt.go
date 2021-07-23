package admin

// StmtInsertEmailVerifier creates a record for email verification request.
// One email only has one record in db.
// If the request's email already exists, perform updates.
// Uniqueness is required since we only want to keep one
// token valid.
const StmtInsertEmailVerifier = `
INSERT INTO b2b.email_verification
SET token = UNHEX(:token),
	email = :email,
	expire_in_days = :expire_in_days,
	created_utc = :created_utc`

const StmtRetrieveEmailVerifier = `
SELECT LOWER(HEX(token)) AS token,
	email,
	expire_in_days,
	created_utc
FROM b2b.email_verification
WHERE token = UNHEX(?)
LIMIT 1`

// StmtEmailVerified set the email_verified to true.
const StmtEmailVerified = `
UPDATE b2b.admin
	SET verified = TRUE
WHERE id = ?
LIMIT 1`
