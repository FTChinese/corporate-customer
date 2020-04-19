package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// Login verifies user email + password.
func (env Env) Login(in admin.AccountInput) (admin.PassportBearer, error) {
	var at admin.Passport

	err := env.db.Get(&at, stmt.Login, in.Email, in.Password)

	if err != nil {
		return admin.PassportBearer{}, err
	}

	jwtAccount, err := admin.NewPassportBearer(at)
	if err != nil {
		return jwtAccount, err
	}
	return jwtAccount, nil
}

// SignUp creates a new admin account.
func (env Env) SignUp(s admin.AccountInput) error {
	_, err := env.db.NamedExec(stmt.SignUp, s)
	if err != nil {
		logger.WithField("trace", "Env.SignUp").Error(err)
	}

	return nil
}

// LoadPassport retrieves and build PassportBearer
// after signup or upon refreshing passport.
func (env Env) LoadPassport(adminID string) (admin.PassportBearer, error) {

	pp, err := env.PassportByAdminID(adminID)
	if err != nil {
		return admin.PassportBearer{}, err
	}

	return admin.NewPassportBearer(pp)
}
