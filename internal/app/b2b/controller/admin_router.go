package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/adminrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"go.uber.org/zap"
)

// AdminRouter defines what an admin can do before login.
type AdminRouter struct {
	guard  B2BJWTGuard
	repo   adminrepo.Env
	post   postman.Postman
	logger *zap.Logger
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(dbs db.ReadWriteMyDBs, p postman.Postman, dk B2BJWTGuard, logger *zap.Logger) AdminRouter {
	return AdminRouter{
		guard:  dk,
		repo:   adminrepo.NewEnv(dbs, logger),
		post:   p,
		logger: logger,
	}
}
