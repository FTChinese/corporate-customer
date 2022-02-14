package cmsrepo

import (
	"github.com/FTChinese/ftacademy/internal/app/repository"
	"github.com/FTChinese/ftacademy/internal/app/repository/txrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap"
)

type Env struct {
	repository.SharedRepo
	logger *zap.Logger
}

func NewEnv(dbs db.ReadWriteMyDBs, logger *zap.Logger) Env {
	return Env{
		SharedRepo: repository.NewSharedRepo(dbs),
		logger:     logger,
	}
}

func (env Env) beginTx() (txrepo.TxRepo, error) {
	tx, err := env.DBs.Write.Beginx()

	if err != nil {
		return txrepo.TxRepo{}, err
	}

	return txrepo.NewTxRepo(tx), nil
}
