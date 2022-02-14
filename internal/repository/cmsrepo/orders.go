package cmsrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/go-rest"
)

func (env Env) listOrders(w pkg.SQLWhere, page gorest.Pagination) ([]checkout.CMSOrderRow, error) {
	var orders = make([]checkout.CMSOrderRow, 0)

	w = w.AddValues(page.Limit, page.Offset())

	err := env.DBs.Read.Select(
		&orders,
		checkout.BuildStmtListOrdersCMS(w.Clause),
		w.Values...)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (env Env) countOrder(w pkg.SQLWhere) (int64, error) {
	var total int64
	err := env.DBs.Read.Get(
		&total,
		checkout.BuildStmtCountOrder(w.Clause),
		w.Values...)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// ListOrders retrieves a list of checkout.CMSOrderRow.
func (env Env) ListOrders(filter checkout.OrderFilter, page gorest.Pagination) (checkout.CMSOrderList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	where := filter.SQLWhere()
	countCh := make(chan int64)
	listCh := make(chan checkout.CMSOrderList)

	go func() {
		defer close(countCh)
		n, err := env.countOrder(where)
		if err != nil {
			sugar.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)

		orders, err := env.listOrders(where, page)

		listCh <- checkout.CMSOrderList{
			PagedList: pkg.PagedList{
				Total:      0,
				Pagination: gorest.Pagination{},
				Err:        err,
			},
			Data: orders,
		}
	}()

	count, listResult := <-countCh, <-listCh
	if listResult.Err != nil {
		return checkout.CMSOrderList{}, listResult.Err
	}

	return checkout.CMSOrderList{
		PagedList: pkg.PagedList{
			Total:      count,
			Pagination: page,
			Err:        nil,
		},
		Data: listResult.Data,
	}, nil
}

func (env Env) LoadOrder(orderID string) (checkout.Order, error) {

	var ord checkout.Order
	err := env.DBs.Read.Get(
		&ord,
		checkout.BuildStmtOrder(false),
		orderID)

	if err != nil {
		return checkout.Order{}, err
	}

	return ord, nil
}

func (env Env) UpdateOrderStatus(o checkout.Order) error {
	_, err := env.DBs.Write.NamedExec(
		checkout.StmtUpdateOrderStatus,
		o)

	if err != nil {
		return err
	}

	return nil
}
