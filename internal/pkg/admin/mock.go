// +build !production

package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

func MockPassportClaims() PassportClaims {
	return PassportClaims{
		AdminID:        uuid.New().String(),
		TeamID:         null.StringFrom(ids.TeamID()),
		StandardClaims: NewStandardClaims(time.Now().Unix() + 86400*7),
	}
}
