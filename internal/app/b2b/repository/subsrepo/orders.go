package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/go-rest"
)

func (env Env) CreateOrder(schema checkout.OrderInputSchema) error {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return err
	}

	err = tx.CreateOrder(schema.OrderRow)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, v := range schema.ItemRows {
		err := tx.CreateOrderItem(v)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		sugar.Error(err)
		return err
	}

	return nil
}

func (env Env) listOrders(teamID string, page gorest.Pagination) ([]checkout.OrderRow, error) {
	var orders = make([]checkout.OrderRow, 0)

	err := env.dbs.Read.Select(
		&orders,
		checkout.StmtListBaseOrders,
		teamID,
		page.Limit,
		page.Offset())

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (env Env) countOrder(teamID string) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(&total, checkout.StmtCountOrder, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) ListOrders(teamID string, page gorest.Pagination) (checkout.OrderRowList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan checkout.OrderRowList)

	go func() {
		defer close(countCh)
		n, err := env.countOrder(teamID)
		if err != nil {
			sugar.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)

		orders, err := env.listOrders(teamID, page)

		listCh <- checkout.OrderRowList{
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
		return checkout.OrderRowList{}, listResult.Err
	}

	return checkout.OrderRowList{
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
func (env Env) orderDetails(r admin.AccessRight) (checkout.Order, error) {
	var ord checkout.Order
	err := env.dbs.Read.Get(&ord, checkout.StmtOrderDetails, r.RowID, r.TeamID)

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

func (env Env) LoadDetailedOrder(r admin.AccessRight) (checkout.Order, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	orderCh := make(chan orderResult)
	itemsCh := make(chan orderItemsResult)

	go func() {
		defer close(orderCh)

		ord, err := env.orderDetails(r)
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

		items, err := env.orderItems(r.RowID)
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
