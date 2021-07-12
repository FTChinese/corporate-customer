package subs

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
	sq2 "github.com/FTChinese/ftacademy/pkg/sq"
	gorest "github.com/FTChinese/go-rest"
)

func (env Env) CreateOrders(o model2.OrderList) error {
	query := stmt.OrderBuilder.Rows(len(o)).Build()

	_, err := env.dbs.Write.Exec(query, sq2.BuildInsertValues(o)...)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadOrder(id, teamID string) (model2.Order, error) {
	var o model2.Order
	err := env.dbs.Read.Get(&o, stmt.GetOrder, id, teamID)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (env Env) ListOrders(teamID string, page gorest.Pagination) ([]model2.Order, error) {
	var o = make([]model2.Order, 0)

	err := env.dbs.Read.Select(&o, stmt.ListOrder, teamID, page.Limit, page.Offset())

	if err != nil {
		return o, err
	}

	return o, nil
}

func (env Env) AsyncListOrders(teamID string, page gorest.Pagination) <-chan model2.PagedOrders {
	r := make(chan model2.PagedOrders)

	go func() {
		defer close(r)

		orders, err := env.ListOrders(teamID, page)

		r <- model2.PagedOrders{
			Data: orders,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountOrder(teamID string) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(&total, stmt.CountOrder, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountOrder(teamID string) <-chan model2.PagedOrders {
	r := make(chan model2.PagedOrders)

	go func() {
		defer close(r)
		total, err := env.CountOrder(teamID)

		r <- model2.PagedOrders{
			Total: total,
			Err:   err,
		}
	}()

	return r
}
