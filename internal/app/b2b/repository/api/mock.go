// +build !production

package api

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/brianvoe/gofakeit/v5"
	"io/ioutil"
)

func MockNewClient() Client {
	config.MustSetupViper()
	return NewSubsAPIClient(false)
}

func (c Client) MustCreateAssignee() licence.Assignee {
	faker.SeedGoFake()
	resp, err := c.ReaderSignup(input.SignupParams{
		Credentials: input.Credentials{
			Email:    gofakeit.Email(),
			Password: faker.SimplePassword(),
		},
	})
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var a licence.Assignee
	if err := json.Unmarshal(b, &a); err != nil {
		panic(err)
	}

	return a
}
