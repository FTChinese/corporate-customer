package txrepo

import "github.com/FTChinese/ftacademy/internal/pkg/checkout"

func (tx TxRepo) CreateOrder(orderRow checkout.OrderRow) error {
	_, err := tx.NamedExec(checkout.StmtCreateBaseOrder, orderRow)
	if err != nil {
		return err
	}

	return nil
}

func (tx TxRepo) CreateOrderItem(item checkout.OrderItem) error {
	_, err := tx.NamedExec(checkout.StmtCreateOrderItem, item)
	if err != nil {
		return err
	}

	return nil
}
