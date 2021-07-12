package subs

import (
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

func (env Env) beginInvTx() (InvitationTx, error) {
	tx, err := env.dbs.Write.Beginx()

	if err != nil {
		return InvitationTx{}, err
	}

	return InvitationTx{
		tx,
	}, nil
}
