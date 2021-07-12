package adminor

import (
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap"
)

type Env struct {
	DBs    db.ReadWriteMyDBs
	logger *zap.Logger
}

func NewEnv(dbs db.ReadWriteMyDBs, logger *zap.Logger) Env {
	return Env{
		DBs:    dbs,
		logger: logger,
	}
}
