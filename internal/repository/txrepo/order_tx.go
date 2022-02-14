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
// Use this together with CreateOrder to save a complete
// shopping cart.
func (tx TxRepo) SaveCartItem(c checkout.CartItemSchema) error {
	_, err := tx.NamedExec(checkout.StmtInsertCartItem, c)
	if err != nil {
		return err
	}

	return nil
}

// SaveLicenceTxnList inserts a list of LicenceTransaction when
// creating an order.
func (tx TxRepo) SaveLicenceTxnList(list checkout.BulkLicenceTxn) error {
	_, err := tx.Exec(
		checkout.StmtBulkLicenceTxn(len(list)).Build(),
		sq.BuildBulkInsertValues(list)...,
	)

	if err != nil {
		return err
	}

	return nil
}
