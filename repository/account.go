package repository

import "github.com/FTChinese/b2b/repository/stmt"

const RetrieveAdminByID = stmt.AccountBase + `
FROM b2b.admin AS a
WHERE a.id = ?
LIMIT 1`

const RetrieveAdminByEmail = stmt.AccountBase + `
FROM b2b.admin AS a
WHERE email = ?
LIMIT 1`

const UpdateDisplayName = `
UPDATE b2b.admin
SET display_name = :display_name,
	updated_utc = UTC_TIMESTAMP()
WHERE id = ?
LIMIT 1`

const UpdatePassword = `
UPDATE b2b.admin
SET password_sha2 = UNHEX(SHA2(:password, 256))
	updated_utc = UTC_TIMESTAMP()
WHERE id = ?
LIMIT 1`

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
