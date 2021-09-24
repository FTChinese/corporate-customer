// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

type Admin struct {
	admin.Account
}

func NewAdmin() Admin {
	return Admin{
		Account: admin.Account{
			BaseAccount: admin.BaseAccount{
				ID:          uuid.New().String(),
				TeamID:      null.StringFrom(ids.TeamID()),
				Email:       faker.GenEmail(),
				DisplayName: null.String{},
				Active:      true,
				Verified:    false,
			},
			Password:   faker.SimplePassword(),
			CreatedUTC: chrono.TimeNow(),
		},
	}
}

func (a Admin) SetTeamID(id string) Admin {
	a.TeamID = null.NewString(id, id != "")
	return a
}

func (a Admin) SetDisplayName(n string) Admin {
	a.DisplayName = null.NewString(n, n != "")
	return a
}

func (a Admin) PassportClaims() admin.PassportClaims {
	return admin.PassportClaims{
		AdminID:        a.ID,
		TeamID:         a.TeamID,
		StandardClaims: admin.NewStandardClaims(time.Now().Unix() + 86400*7),
	}
}

func (a Admin) Creator() admin.Creator {
	return admin.Creator{
		AdminID: a.ID,
		TeamID:  a.TeamID.String,
	}
}

func (a Admin) EmailVerifier() admin.EmailVerifier {
	v, err := admin.NewEmailVerifier(a.Email)
	if err != nil {
		panic(err)
	}

	return v
}

func (a Admin) PwResetSession() admin.PwResetSession {
	s, err := admin.NewPwResetSession(input.ForgotPasswordParams{
		Email:     a.Email,
		SourceURL: null.String{},
	})
	if err != nil {
		panic(err)
	}

	return s
}

func (a Admin) TeamParams() input.TeamParams {
	faker.SeedGoFake()

	return input.TeamParams{
		OrgName:      gofakeit.Company(),
		InvoiceTitle: null.StringFrom(gofakeit.Company()),
	}
}

func (a Admin) Team() admin.Team {
	t := admin.NewTeam(a.ID, a.TeamParams())
	t.ID = a.TeamID.String

	return t
}
