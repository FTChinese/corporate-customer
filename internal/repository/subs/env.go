package subs

import (
	txrepo2 "github.com/FTChinese/ftacademy/internal/repository/txrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap"
)

type Env struct {
	dbs    db.ReadWriteMyDBs
	logger *zap.Logger
}

func NewEnv(DBs db.ReadWriteMyDBs, logger *zap.Logger) Env {
	return Env{
		dbs:    DBs,
		logger: logger,
	}
}

func (env Env) beginTx() (txrepo2.TxRepo, error) {
	tx, err := env.dbs.Write.Beginx()

	if err != nil {
		return txrepo2.TxRepo{}, err
	}

	return txrepo2.NewTxRepo(tx), nil
}
