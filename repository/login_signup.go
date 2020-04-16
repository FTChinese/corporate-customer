package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// Login verifies user email + password.
// TODO: cache AccountTeam by admin id.
func (env Env) Login(in admin.AccountInput) (admin.JWTAccount, error) {
	var at admin.AccountTeam

	err := env.db.Get(&at, stmt.Login, in.Email, in.Password)

	if err != nil {
		return admin.JWTAccount{}, err
	}

	jwtAccount, err := admin.NewJWTAccount(at)
	if err != nil {
		return jwtAccount, err
	}
	return jwtAccount, nil
}

func (env Env) SignUp(s admin.AccountInput) error {
	_, err := env.db.NamedExec(stmt.SignUp, s)
	if err != nil {
		logger.WithField("trace", "Env.SignUp").Error(err)
	}

	return nil
}

// JWTAccount retrieves and build JWTAccount used to refresh
// jwt, or signup.
func (env Env) JWTAccount(id string) (admin.JWTAccount, error) {

	at, err := env.AccountTeam(id)
	if err != nil {
		return admin.JWTAccount{}, err
	}

	return admin.NewJWTAccount(at)
}
