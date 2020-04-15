package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/sq"
	"github.com/FTChinese/b2b/repository/stmt"
)

func (env Env) CreateOrders(o admin.OrderList) error {
	query := stmt.OrderBuilder.Rows(len(o)).Build()

	_, err := env.db.Exec(query, sq.BuildInsertValues(o)...)
	if err != nil {
		return err
	}

	return nil
}

const stmtOrder = stmt.Order + `
WHERE id = ?
	AND team_id = ?
LIMIT 1`

func (env Env) LoadOrder(id, teamID string) (admin.Order, error) {
	var o admin.Order
	err := env.db.Get(&o, stmtOrder, id, teamID)

	if err != nil {
		return o, err
	}

	return o, nil
}

const stmtListOrder = stmt.Order + `
WHERE id = ?
ORDER BY created_utc DESC`

func (env Env) ListOrders(teamID string) ([]admin.Order, error) {
	var o = make([]admin.Order, 0)

	err := env.db.Select(&o, teamID)

	if err != nil {
		return o, err
	}

	return o, nil
}
