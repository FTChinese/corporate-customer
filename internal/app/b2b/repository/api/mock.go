//go:build !production
// +build !production

package api

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/brianvoe/gofakeit/v5"
	"net/http"
)

func MockNewClient() Client {
	faker.MustSetupViper()
	return NewSubsAPIClient(false)
}

func (c Client) MustCreateAssignee() licence.Assignee {
	faker.SeedGoFake()

	resp, err := c.EmailSignUp(faker.MustMarshalIndent(input.SignupParams{
		Credentials: input.Credentials{
			Email:    gofakeit.Email(),
			Password: faker.SimplePassword(),
		},
	}), http.Header{})

	if err != nil {
		panic(err)
	}

	var a licence.Assignee
	if err := json.Unmarshal(resp.Body, &a); err != nil {
		panic(err)
	}

	return a
}
