package repository

import "github.com/FTChinese/b2b/models/admin"

const stmtSignUp = `
INSERT INTO b2b.admin
SET id = :admin_id,
	email = :email,
	password_sha2 = UNHEX(SHA2(:password, 256)),
	vrf_token = UNHEX(:token),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

func (env Env) SignUp(s admin.SignUp) error {
	_, err := env.db.NamedExec(stmtSignUp, s)
	if err != nil {
		logger.WithField("trace", "Env.SignUp").Error(err)
	}

	return nil
}

// LogIn verifies email and password combination.
const stmtLogIn = `
SELECT id AS admin_id
FROM b2b.admin
WHERE (email, password_sha2) = (?, UNHEX(SHA2(?, 256)))
LIMIT 1`

// Login verifies user password and returns id.
func (env Env) Login(l admin.Login) (string, error) {
	var id string

	err := env.db.Get(&id, stmtLogIn, l.Email, l.Password)

	if err != nil {
		return id, err
	}

	return id, nil
}
