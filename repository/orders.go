package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

func (env Env) CreateOrder(o admin.Order) error {
	_, err := env.db.NamedExec(stmt.CreateOrder, o)
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
