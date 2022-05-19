package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
)

func (c Client) ListLegalDoc(rawQuery string) (pkg.LegalList, *render.ResponseError) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(basePathLegal).
		SetRawQuery(rawQuery).
		String()

	var l pkg.LegalList
	respErr := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		EndJSON(&l)

	if respErr != nil {
		return pkg.LegalList{}, respErr
	}

	return l, nil
}

func (c Client) LoadLegalDoc(id string) (pkg.LegalDoc, *render.ResponseError) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(basePathLegal).
		AddPath(id).
		String()

	var d pkg.LegalDoc
	respErr := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		EndJSON(&d)

	if respErr != nil {
		return pkg.LegalDoc{}, respErr
	}

	return d, nil

}
