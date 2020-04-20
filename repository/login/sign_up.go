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

// LoadPassport retrieves and build PassportBearer
// after signup.
func (env Env) LoadPassport(adminID string) (admin.PassportBearer, error) {

	var pp admin.Passport
	if err := env.db.Get(&pp, stmt.PassportByAdminID, adminID); err != nil {
		return admin.PassportBearer{}, err
	}

	return admin.NewPassportBearer(pp)
}
