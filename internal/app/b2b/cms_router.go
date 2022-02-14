package b2b

import (
	"github.com/FTChinese/ftacademy/internal/repository/cmsrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"go.uber.org/zap"
)

type CMSRouter struct {
	repo   cmsrepo.Env
	post   postman.Postman
	logger *zap.Logger
}

func NewCMSRouter(dbs db.ReadWriteMyDBs, p postman.Postman, logger *zap.Logger) CMSRouter {
	return CMSRouter{
		repo:   cmsrepo.NewEnv(dbs, logger),
		post:   p,
		logger: logger,
	}
}
