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
	return NewClients(false).Select(false)
}

func (c Client) MustCreateAssignee() licence.Assignee {
	faker.SeedGoFake()

	resp, err := c.EmailSignUp(faker.MustMarshalToReader(input.SignupParams{
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
