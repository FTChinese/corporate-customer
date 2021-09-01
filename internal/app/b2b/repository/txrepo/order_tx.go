package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/pkg/sq"
)

// CreateOrder saves a row into order table.
func (tx TxRepo) CreateOrder(order checkout.Order) error {
	_, err := tx.NamedExec(checkout.StmtCreateOrder, order)
	if err != nil {
		return err
	}

	return nil
}

// SaveCartItem saves an element of shopping cart's
// item array.
// Used this together with CreateOrder to save a complete
// shopping cart.
func (tx TxRepo) SaveCartItem(c checkout.CartItemSchema) error {
	_, err := tx.NamedExec(checkout.StmtInsertCartItem, c)
	if err != nil {
		return err
	}

	return nil
}

// SaveLicenceQueue inserts a list of LicenceQueue when
// creating an order.
func (tx TxRepo) SaveLicenceQueue(q checkout.BulkLicenceQueue) error {
	_, err := tx.Exec(
		checkout.StmtBulkLicenceQueue(len(q)).Build(),
		sq.BuildBulkInsertValues(q),
	)

	if err != nil {
		return err
	}

	return nil
}
