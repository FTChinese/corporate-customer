package txrepo

import "github.com/FTChinese/ftacademy/internal/pkg/checkout"

func (tx TxRepo) CreateOrder(bo checkout.BaseOrder) error {
	_, err := tx.NamedExec(checkout.StmtCreateBaseOrder, bo)
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
