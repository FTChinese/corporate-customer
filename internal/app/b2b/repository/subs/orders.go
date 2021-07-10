package subs

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	sq2 "github.com/FTChinese/ftacademy/pkg/sq"
	gorest "github.com/FTChinese/go-rest"
)

func (env Env) CreateOrders(o model.OrderList) error {
	query := stmt.OrderBuilder.Rows(len(o)).Build()

	_, err := env.db.Exec(query, sq2.BuildInsertValues(o)...)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadOrder(id, teamID string) (model.Order, error) {
	var o model.Order
	err := env.db.Get(&o, stmt.GetOrder, id, teamID)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (env Env) ListOrders(teamID string, page gorest.Pagination) ([]model.Order, error) {
	var o = make([]model.Order, 0)

	err := env.db.Select(&o, stmt.ListOrder, teamID, page.Limit, page.Offset())

	if err != nil {
		logger.WithField("trace", "ListOrders").Error(err)
		return o, err
	}

	return o, nil
}

func (env Env) AsyncListOrders(teamID string, page gorest.Pagination) <-chan model.PagedOrders {
	r := make(chan model.PagedOrders)

	go func() {
		defer close(r)

		orders, err := env.ListOrders(teamID, page)

		r <- model.PagedOrders{
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

func (env Env) AsyncCountOrder(teamID string) <-chan model.PagedOrders {
	r := make(chan model.PagedOrders)

	go func() {
		defer close(r)
		total, err := env.CountOrder(teamID)

		r <- model.PagedOrders{
			Total: total,
			Err:   err,
		}
	}()

	return r
}
