package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/adminrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"go.uber.org/zap"
)

type CMSRouter struct {
	adminRepo adminrepo.Env
	post      postman.Postman
	logger    *zap.Logger
}

func NewCMSRouter(dbs db.ReadWriteMyDBs, p postman.Postman, logger *zap.Logger) CMSRouter {
	return CMSRouter{
		adminRepo: adminrepo.NewEnv(dbs, logger),
		post:      p,
		logger:    logger,
	}
}
