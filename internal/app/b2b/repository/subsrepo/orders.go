package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/go-rest"
)

// CreateOrder inserts a row into order table based on
// shopping cart.
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
		sugar.Error(err)
		_ = tx.Rollback()
		return err
	}

	for _, v := range schema.CartItems {
		err := tx.SaveCartItem(v)
		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.SaveLicenceQueue(schema.Queue)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		sugar.Error(err)
		return err
	}

	return nil
}

// Retrieve a list of orders belong to a team.
func (env Env) listOrders(teamID string, page gorest.Pagination) ([]checkout.Order, error) {
	var orders = make([]checkout.Order, 0)

	err := env.DBs.Read.Select(
		&orders,
		checkout.StmtListOrders,
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
	err := env.DBs.Read.Get(
		&total,
		checkout.StmtCountOrder,
		teamID)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// ListOrders retrieves a list of orders created by an admin.
func (env Env) ListOrders(teamID string, page gorest.Pagination) (checkout.OrderList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan checkout.OrderList)

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

		listCh <- checkout.OrderList{
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
		return checkout.OrderList{}, listResult.Err
	}

	return checkout.OrderList{
		PagedList: pkg.PagedList{
			Total:      count,
			Pagination: page,
			Err:        nil,
		},
		Data: listResult.Data,
	}, nil
}

func (env Env) RetrieveOrder(r admin.AccessRight) (checkout.Order, error) {

	var ord checkout.Order
	err := env.DBs.Read.Get(
		&ord,
		checkout.BuildStmtOrder(true),
		r.RowID,
		r.TeamID,
	)

	if err != nil {
		return checkout.Order{}, err
	}

	return ord, nil
}

func (env Env) RetrieveGroupedQueue(orderID string, priceID string) (checkout.GroupedQueues, error) {
	var q = make([]checkout.LicenceQueue, 0)
	err := env.DBs.Read.Select(
		&q,
		checkout.StmtListLicenceQueue,
		orderID,
		priceID,
	)
	if err != nil {
		return checkout.GroupedQueues{}, err
	}

	return checkout.NewGroupedQueues(priceID, q), nil
}
