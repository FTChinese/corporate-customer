package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/adminor"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"go.uber.org/zap"
)

// AdminRouter defines what an admin can do before login.
type AdminRouter struct {
	keeper Doorkeeper
	repo   adminor.Env
	post   postman.Postman
	logger *zap.Logger
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(dbs db.ReadWriteMyDBs, p postman.Postman, dk Doorkeeper, logger *zap.Logger) AdminRouter {
	return AdminRouter{
		keeper: dk,
		repo:   adminor.NewEnv(dbs, logger),
		post:   p,
		logger: logger,
	}
}
