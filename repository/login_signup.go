package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// Login verifies user email + password.
func (env Env) Login(in admin.AccountInput) (admin.JWTAccount, error) {
	var a admin.JWTAccount

	err := env.db.Get(&a, stmt.Login, in.Email, in.Password)

	if err != nil {
		return admin.JWTAccount{}, err
	}

	a, err = a.WithToken()
	if err != nil {
		return a, err
	}
	return a, nil
}

const stmtSignUp = `
INSERT INTO b2b.admin
SET id = :admin_id,
	email = :email,
	password_sha2 = UNHEX(SHA2(:password, 256)),
	vrf_token = UNHEX(:token),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

func (env Env) SignUp(s admin.AccountInput) error {
	_, err := env.db.NamedExec(stmtSignUp, s)
	if err != nil {
		logger.WithField("trace", "Env.SignUp").Error(err)
	}

	return nil
}

func (env Env) JWTAccount(id string) (admin.JWTAccount, error) {
	var a admin.JWTAccount
	if err := env.db.Get(&a, stmt.JWTAccount, id); err != nil {
		return a, err
	}

	a, err := a.WithToken()
	if err != nil {
		return a, err
	}

	return a, nil
}
