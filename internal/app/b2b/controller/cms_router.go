package controller

import (
	adminor2 "github.com/FTChinese/ftacademy/internal/repository/adminor"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"go.uber.org/zap"
)

type CMSRouter struct {
	adminRepo adminor2.Env
	post      postman.Postman
	logger    *zap.Logger
}

func NewCMSRouter(dbs db.ReadWriteMyDBs, p postman.Postman, logger *zap.Logger) CMSRouter {
	return CMSRouter{
		adminRepo: adminor2.NewEnv(dbs, logger),
		post:      p,
		logger:    logger,
	}
}
