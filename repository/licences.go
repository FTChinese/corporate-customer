package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

const stmtCreateLicence = `
INSERT INTO b2b.licence
SET id = :licence_id
	team_id = :team_id,
	plan_id = :plan_id,
	created_utc = UTC_TIMESTAMP()`

func (env Env) CreateLicence(l admin.Licence) error {
	_, err := env.db.NamedExec(stmtCreateLicence, l)
	if err != nil {
		return err
	}

	return nil
}

// LoadLicence retrieves an expanded licence
// belong to a team.
func (env Env) LoadLicence(id, teamID string) (admin.ExpandedLicence, error) {
	var ls admin.LicenceSchema
	err := env.db.Get(&ls, stmt.ExpandedLicence, id, teamID)

	if err != nil {
		return admin.ExpandedLicence{}, err
	}

	return ls.ExpandedLicence(), nil
}

func (env Env) ListLicence(teamID string) ([]admin.ExpandedLicence, error) {
	var ls = make([]admin.LicenceSchema, 0)

	err := env.db.Select(&ls, stmt.ListExpandedLicences, teamID)

	if err != nil {
		return nil, err
	}

	el := make([]admin.ExpandedLicence, 0)
	for _, v := range ls {
		el = append(el, v.ExpandedLicence())
	}
	return el, nil
}
