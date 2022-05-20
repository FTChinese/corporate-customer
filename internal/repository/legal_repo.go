package repository

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/go-rest/render"
	"github.com/patrickmn/go-cache"
)

type LegalRepo struct {
	client api.Client
	CacheRepo
}

func NewLegalRepo(client api.Client, c *cache.Cache) LegalRepo {
	return LegalRepo{
		client:    client,
		CacheRepo: NewCacheRepo(c),
	}
}

func (repo LegalRepo) ListLegalDoc(rawQuery string) (pkg.LegalList, *render.ResponseError) {
	return repo.client.ListLegalDoc(rawQuery)
}

func (repo LegalRepo) LoadOrFetchLegalDoc(id string, refresh bool) (pkg.LegalDoc, *render.ResponseError) {
	if !refresh {
		doc, err := repo.LoadLegalDoc(id)
		if err == nil {
			return doc, nil
		}
	}

	doc, err := repo.client.LoadLegalDoc(id)
	if err != nil {
		return pkg.LegalDoc{}, err
	}

	rendered := doc.Rendered()

	repo.SaveLegalDoc(rendered)

	return rendered, nil
}
