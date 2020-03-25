package stmt

const retrieveAdmin = `
SELECT a.id AS admin_id,
	a.email,
	a.display_name`

const RetrieveAdminByID = retrieveAdmin + `
FROM b2b.admin AS a
WHERE id = ?
LIMIT 1`

const RetrieveAdminByEmail = retrieveAdmin + `
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
