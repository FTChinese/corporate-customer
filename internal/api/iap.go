package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"net/http"
)

func (c Client) RefreshIAP(ids reader.PassportClaims, originalTxID string) (*http.Response, error) {
	url := c.baseURL + pathAppleSubOf(originalTxID)

	resp, errs := fetch.
		New().
		Patch(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
