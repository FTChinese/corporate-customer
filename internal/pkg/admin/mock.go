// +build !production

package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

func MockPassportClaims() PassportClaims {
	return PassportClaims{
		AdminID:        uuid.New().String(),
		TeamID:         null.StringFrom(pkg.TeamID()),
		StandardClaims: NewStandardClaims(time.Now().Unix() + 86400*7),
	}
}
