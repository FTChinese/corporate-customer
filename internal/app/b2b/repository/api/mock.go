// +build !production

package api

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/config"
	"io/ioutil"
)

func MockNewClient() Client {
	config.MustSetupViper()
	return NewSubsAPIClient(false)
}

func (c Client) MustCreateAssignee(s input.SignupParams) licence.Assignee {
	resp, err := c.ReaderSignup(s)
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
