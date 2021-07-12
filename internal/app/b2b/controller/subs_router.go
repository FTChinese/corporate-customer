package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"go.uber.org/zap"
)

type SubsRouter struct {
	repo   subs.Env
	post   postman.Postman
	logger *zap.Logger
}

func NewSubsRouter(myDBs db.ReadWriteMyDBs, pm postman.Postman, logger *zap.Logger) SubsRouter {
	return SubsRouter{
		repo:   subs.NewEnv(myDBs, logger),
		post:   pm,
		logger: logger,
	}
}
