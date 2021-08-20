package cmsrepo

import (
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap"
)

type Env struct {
	dbs    db.ReadWriteMyDBs
	logger *zap.Logger
}

func NewEnv(dbs db.ReadWriteMyDBs, logger *zap.Logger) Env {
	return Env{
		dbs:    dbs,
		logger: logger,
	}
}
