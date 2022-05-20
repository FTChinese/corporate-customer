package content

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/internal/repository"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"time"
)

type Routes struct {
	repo    repository.LegalRepo
	version string
	logger  *zap.Logger
}

func NewRoutes(client api.Client, version string, logger *zap.Logger) Routes {
	return Routes{
		repo:    repository.NewLegalRepo(client, cache.New(24*time.Hour, 1*time.Hour)),
		version: version,
		logger:  logger,
	}
}
