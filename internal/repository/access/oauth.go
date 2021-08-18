package access

import (
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/oauth"
	"github.com/patrickmn/go-cache"
	"time"
)

type Env struct {
	dbs   db.ReadWriteMyDBs
	cache *cache.Cache
}

func NewEnv(dbs db.ReadWriteMyDBs) Env {
	return Env{
		dbs: dbs,
		// Default expiration 24 hours, and purges the expired items every hour.
		cache: cache.New(24*time.Hour, 1*time.Hour),
	}
}

// Load tries to load an access token from cache first, then
// retrieve from db if not found in cache.
func (env Env) Load(token string) (oauth.OAuth, error) {
	if acc, ok := env.loadCachedToken(token); ok {
		return acc, nil
	}

	acc, err := env.retrieveFromDB(token)
	if err != nil {
		return acc, err
	}

	env.cacheToken(token, acc)

	return acc, nil
}

func (env Env) loadCachedToken(token string) (oauth.OAuth, bool) {
	x, found := env.cache.Get(token)
	if !found {
		return oauth.OAuth{}, false
	}

	if access, ok := x.(oauth.OAuth); ok {
		return access, true
	}

	return oauth.OAuth{}, false
}

func (env Env) retrieveFromDB(token string) (oauth.OAuth, error) {
	var access oauth.OAuth

	if err := env.dbs.Read.Get(&access, oauth.StmtOAuth, token); err != nil {
		return access, err
	}

	return access, nil
}

func (env Env) cacheToken(token string, access oauth.OAuth) {
	env.cache.Set(token, access, cache.DefaultExpiration)
}
