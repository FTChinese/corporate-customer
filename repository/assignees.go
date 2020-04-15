package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/reader"
	"github.com/FTChinese/b2b/repository/stmt"
)

func (env Env) FindReader(email string) (reader.Reader, error) {
	var r reader.Reader
	err := env.db.Get(&r, stmt.SelectReader, email)
	if err != nil {
		return r, err
	}

	r.Normalize()

	return r, nil
}

const stmtSaveAssignee = `
INSERT INTO b2b.assignee
SET email = :email,
	ftc_id = :ftc_id,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

func (env Env) SaveAssignee(a admin.AssigneeSchema) error {
	_, err := env.db.NamedExec(stmtSaveAssignee, a)

	if err != nil {
		return err
	}

	return nil
}
