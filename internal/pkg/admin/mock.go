// +build !production

package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

func MockAccount() Account {
	return Account{
		BaseAccount: BaseAccount{
			ID:          uuid.New().String(),
			TeamID:      null.StringFrom(pkg.TeamID()),
			Email:       faker.GenEmail(),
			DisplayName: null.StringFrom(gofakeit.Username()),
			Active:      false,
			Verified:    false,
		},
		Password:   faker.SimplePassword(),
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}
