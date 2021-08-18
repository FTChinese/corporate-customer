package oauth

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"time"
)

const StmtOAuth = `
SELECT access_token,
	is_active,
	expires_in,
	created_utc
FROM oauth.access
WHERE access_token = UNHEX(?)
LIMIT 1`

// OAuth contains the data related to an oauth access token,
// used either by human or machines.
type OAuth struct {
	Token     string      `db:"access_token"`
	Active    bool        `db:"is_active"`
	ExpiresIn null.Int    `db:"expires_in"` // seconds
	CreatedAt chrono.Time `db:"created_utc"`
}

func (a OAuth) Expired() bool {

	if a.ExpiresIn.IsZero() {
		return false
	}

	expireAt := a.CreatedAt.Add(time.Second * time.Duration(a.ExpiresIn.Int64))

	if expireAt.Before(time.Now()) {
		return true
	}

	return false
}
