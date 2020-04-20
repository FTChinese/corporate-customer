package login

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// Login verifies user email + password.
func (env Env) Login(in admin.AccountInput) (admin.PassportBearer, error) {
	var at admin.Passport

	err := env.db.Get(&at, stmt.Login, in.Email, in.Password)

	if err != nil {
		logger.WithField("trace", "Login").Error(err)
		return admin.PassportBearer{}, err
	}

	jwtAccount, err := admin.NewPassportBearer(at)
	if err != nil {
		logger.WithField("trace", "Login").Error(err)
		return jwtAccount, err
	}
	return jwtAccount, nil
}
