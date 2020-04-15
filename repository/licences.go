package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// LoadLicence retrieves a licence, together with its
// subscription plan and the user to whom it was assigned.
// If the licence is not assigned yet, user part should
// have empty value.
func (env Env) LoadLicence(id, teamID string) (admin.ExpandedLicence, error) {
	var ls admin.LicenceSchema
	err := env.db.Get(&ls, stmt.ExpandedLicence, id, teamID)

	if err != nil {
		return admin.ExpandedLicence{}, err
	}

	return ls.ExpandedLicence()
}

func (env Env) ListLicence(teamID string) ([]admin.ExpandedLicence, error) {
	var ls = make([]admin.LicenceSchema, 0)

	err := env.db.Select(&ls, stmt.ListExpandedLicences, teamID)

	if err != nil {
		return nil, err
	}

	el := make([]admin.ExpandedLicence, 0)
	for _, v := range ls {
		l, err := v.ExpandedLicence()
		if err != nil {
			return nil, err
		}
		el = append(el, l)
	}
	return el, nil
}
