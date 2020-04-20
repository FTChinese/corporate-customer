package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/sq"
	"github.com/FTChinese/b2b/repository/stmt"
	gorest "github.com/FTChinese/go-rest"
)

func (env Env) CreateOrders(o admin.OrderList) error {
	query := stmt.OrderBuilder.Rows(len(o)).Build()

	_, err := env.db.Exec(query, sq.BuildInsertValues(o)...)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadOrder(id, teamID string) (admin.Order, error) {
	var o admin.Order
	err := env.db.Get(&o, stmt.GetOrder, id, teamID)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (env Env) ListOrders(teamID string, page gorest.Pagination) ([]admin.Order, error) {
	var o = make([]admin.Order, 0)

	err := env.db.Select(&o, stmt.ListOrder, teamID, page.Limit, page.Offset())

	if err != nil {
		logger.WithField("trace", "ListOrders").Error(err)
		return o, err
	}

	return o, nil
}

func (env Env) AsyncListOrders(teamID string, page gorest.Pagination) <-chan admin.PagedOrders {
	r := make(chan admin.PagedOrders)

	go func() {
		defer close(r)

		orders, err := env.ListOrders(teamID, page)

		r <- admin.PagedOrders{
			Data: orders,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountOrder(teamID string) (int64, error) {
	var total int64
	err := env.db.Get(&total, stmt.CountOrder, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountOrder(teamID string) <-chan admin.PagedOrders {
	r := make(chan admin.PagedOrders)

	go func() {
		defer close(r)
		total, err := env.CountOrder(teamID)

		r <- admin.PagedOrders{
			Total: total,
			Err:   err,
		}
	}()

	return r
}
