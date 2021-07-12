package admin

const colAuthResult = `
SELECT id,
	password = UNHEX(SHA2(?, 256)) AS password_matched
FROM b2b.admin
`

const StmtVerifyPasswordByEmail = colAuthResult + `
WHERE email = ?
LIMIT 1`

const StmtVerifyPasswordByID = colAuthResult + `
WHERE id = ?
LIMIT 1`
