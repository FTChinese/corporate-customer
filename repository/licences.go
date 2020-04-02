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

const stmtLicence = stmt.Licence + `
WHERE id = ?
	AND team_id = ?
LIMIT 1`

func (env Env) LoadLicence(id, teamID string) (admin.Licence, error) {
	var l admin.Licence
	err := env.db.Get(&l, stmtLicence, id, teamID)

	if err != nil {
		return l, err
	}

	return l, nil
}

const stmtListLicences = stmt.Licence + `
WHERE team_id = ?`

func (env Env) ListLicence(teamID string) ([]admin.Licence, error) {
	var l = make([]admin.Licence, 0)

	err := env.db.Select(&l, stmtLicence, teamID)

	if err != nil {
		return l, err
	}

	return l, nil
}
