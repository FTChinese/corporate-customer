package content

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/internal/repository"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/FTChinese/ftacademy/web"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"net/http"
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

func (routes Routes) ListLegalDoc(c echo.Context) error {
	query := c.QueryString()
	builder := web.NewContextBuilder(routes.version).SetTitle("条款")
	docList, err := routes.repo.ListLegalDoc(query)

	if err != nil {
		builder.SetErr(err.Message)
	} else {
		builder.SetLegalList(docList.Data)
	}

	return c.Render(http.StatusOK, "legal/home.html", builder.Build())
}

func (routes Routes) LoadLegalDoc(c echo.Context) error {
	refresh := xhttp.GetQueryRefresh(c)
	id := c.Param("id")
	builder := web.NewContextBuilder(routes.version)
	doc, err := routes.repo.LoadOrFetchLegalDoc(id, refresh)

	if err != nil {
		builder.SetErr(err.Message)
	} else {
		builder.SetLegalDoc(doc).SetTitle(doc.Title)
	}

	return c.Render(http.StatusOK, "legal/legal_doc.html", builder.Build())
}
