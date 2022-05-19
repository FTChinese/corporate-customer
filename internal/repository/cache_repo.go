package repository

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/patrickmn/go-cache"
)

type CacheRepo struct {
	cache *cache.Cache
}

func NewCacheRepo(c *cache.Cache) CacheRepo {
	return CacheRepo{cache: c}
}

func (repo CacheRepo) SaveLegalDoc(d pkg.LegalDoc) {
	repo.cache.Set(d.ID, d, cache.DefaultExpiration)
}

func (repo CacheRepo) LoadLegalDoc(id string) (pkg.LegalDoc, error) {
	x, found := repo.cache.Get(id)
	if found {
		if d, ok := x.(pkg.LegalDoc); ok {
			return d, nil
		}
	}

	return pkg.LegalDoc{}, errors.New("not found")
}

func (repo CacheRepo) Remove(key string) {
	repo.cache.Delete(key)
}
