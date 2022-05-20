package content

import (
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/FTChinese/ftacademy/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

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

// LoadLegalDoc loads a legal document from cache, and fallback to API if not found.
// It will be cached in memory if fetched from API.
// When optional query parameter `?refresh=true` is passed in,
// cached version of this document will be dropped and a new request to API will be sent.
// Use the refresh feature upon document editing.
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
