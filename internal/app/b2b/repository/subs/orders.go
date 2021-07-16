package subs

import (
	"github.com/FTChinese/ftacademy/internal/pkg/order"
	gorest "github.com/FTChinese/go-rest"
)

func (env Env) LoadOrder(id, teamID string) (order.Order, error) {
	var o order.Order
	err := env.dbs.Read.Get(&o, order.GetOrder, id, teamID)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (env Env) ListOrders(teamID string, page gorest.Pagination) ([]order.Order, error) {
	var o = make([]order.Order, 0)

	err := env.dbs.Read.Select(&o, order.ListOrder, teamID, page.Limit, page.Offset())

	if err != nil {
		return o, err
	}

	return o, nil
}

func (env Env) AsyncListOrders(teamID string, page gorest.Pagination) <-chan order.PagedOrders {
	r := make(chan order.PagedOrders)

	go func() {
		defer close(r)

		orders, err := env.ListOrders(teamID, page)

		r <- order.PagedOrders{
			Data: orders,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountOrder(teamID string) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(&total, order.CountOrder, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountOrder(teamID string) <-chan order.PagedOrders {
	r := make(chan order.PagedOrders)

	go func() {
		defer close(r)
		total, err := env.CountOrder(teamID)

		r <- order.PagedOrders{
			Total: total,
			Err:   err,
		}
	}()

	return r
}
