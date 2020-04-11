package repository

import "github.com/FTChinese/b2b/models/admin"

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
