// +build !production

package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
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
			TeamID:      null.String{},
			Email:       faker.GenEmail(),
			DisplayName: null.String{},
			Active:      true,
			Verified:    false,
		},
		Password:   faker.SimplePassword(),
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}

func MockEmailVerifier(email string) EmailVerifier {
	v, err := NewEmailVerifier(email)
	if err != nil {
		panic(err)
	}

	return v
}

func MockPwResetSession(params input.ForgotPasswordParams) PwResetSession {
	s, err := NewPwResetSession(params)
	if err != nil {
		panic(err)
	}

	return s
}

func MockTeamParams() input.TeamParams {
	faker.SeedGoFake()

	return input.TeamParams{
		OrgName:      gofakeit.Company(),
		InvoiceTitle: null.StringFrom(gofakeit.Company()),
	}
}

func MockTeam() Team {
	return NewTeam(uuid.New().String(), MockTeamParams())
}
