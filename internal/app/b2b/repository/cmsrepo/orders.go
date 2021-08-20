package cmsrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/go-rest"
)

func (env Env) listOrders(w pkg.SQLWhere, page gorest.Pagination) ([]checkout.CMSOrderRow, error) {
	var orders = make([]checkout.CMSOrderRow, 0)

	w = w.AddValues(page.Limit, page.Offset())

	err := env.dbs.Read.Select(
		&orders,
		checkout.BuildStmtListOrders(w.Clause),
		w.Values...)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (env Env) countOrder(w pkg.SQLWhere) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(
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

// orderDetails retrieve a row from order table, excluding
// checkout_products column.
func (env Env) orderDetails(orderID string) (checkout.Order, error) {
	var ord checkout.Order
	err := env.dbs.Read.Get(
		&ord,
		checkout.BuildStmtOrder(false),
		orderID)

	if err != nil {
		return checkout.Order{}, err
	}

	return ord, nil
}

func (env Env) orderItems(orderID string) ([]checkout.OrderItem, error) {
	var items = make([]checkout.OrderItem, 0)
	err := env.dbs.Read.Select(
		&items,
		checkout.StmtItemsOfOrder,
		orderID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type orderResult struct {
	value checkout.Order
	err   error
}

type orderItemsResult struct {
	value []checkout.OrderItem
	err   error
}

func (env Env) LoadDetailedOrder(orderID string) (checkout.Order, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	orderCh := make(chan orderResult)
	itemsCh := make(chan orderItemsResult)

	go func() {
		defer close(orderCh)

		ord, err := env.orderDetails(orderID)
		if err != nil {
			sugar.Error(err)
		}
		orderCh <- orderResult{
			value: ord,
			err:   err,
		}
	}()

	go func() {
		defer close(itemsCh)

		items, err := env.orderItems(orderID)
		if err != nil {
			sugar.Error(err)
		}
		itemsCh <- orderItemsResult{
			value: items,
			err:   err,
		}
	}()

	ordRes, itemsRes := <-orderCh, <-itemsCh
	if ordRes.err != nil {
		return checkout.Order{}, ordRes.err
	}
	if itemsRes.err != nil {
		return checkout.Order{}, itemsRes.err
	}

	return checkout.Order{
		BaseOrder: ordRes.value.BaseOrder,
		Items:     itemsRes.value,
		Payment:   ordRes.value.Payment,
	}, nil
}
