package txrepo

import (
	"github.com/jmoiron/sqlx"
)

type TxRepo struct {
	*sqlx.Tx
}

func NewTxRepo(tx *sqlx.Tx) TxRepo {
	return TxRepo{
		Tx: tx,
	}
}
