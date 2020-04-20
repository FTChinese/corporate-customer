package repository

import (
	"github.com/FTChinese/b2b/models/admin"
)

// LoadPassport retrieves and build PassportBearer
// after signup or upon refreshing passport.
func (env Env) LoadPassport(adminID string) (admin.PassportBearer, error) {

	pp, err := env.PassportByAdminID(adminID)
	if err != nil {
		return admin.PassportBearer{}, err
	}

	return admin.NewPassportBearer(pp)
}
