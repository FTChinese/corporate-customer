package login

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// SignUp creates a new admin account.
func (env Env) SignUp(s admin.AccountInput) error {
	_, err := env.db.NamedExec(stmt.SignUp, s)
	if err != nil {
		logger.WithField("trace", "SignUp").Error(err)
	}

	return nil
}

// PassportByAdminID retrieves admin's account and team data
// by admin id.
func (env Env) adminPassport(id string) (admin.Passport, error) {
	var a admin.Passport
	if err := env.db.Get(&a, stmt.PassportByAdminID, id); err != nil {
		return a, err
	}

	return a, nil
}

// LoadPassport retrieves and build PassportBearer
// after signup.
func (env Env) LoadPassport(adminID string) (admin.PassportBearer, error) {

	pp, err := env.adminPassport(adminID)
	if err != nil {
		return admin.PassportBearer{}, err
	}

	return admin.NewPassportBearer(pp)
}
