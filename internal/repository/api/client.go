package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"log"
	"net/http"
)

type Client struct {
	key     string
	baseURL string
}

func NewSubsAPIClient(prod bool) Client {
	return Client{
		key:     config.MustSubsAPIKey().Pick(prod),
		baseURL: config.MustSubsAPIv3BaseURL().Pick(prod),
	}
}

func (c Client) ReaderSignup(s input.SignupParams) (*http.Response, error) {
	url := c.baseURL + "/auth/email/signup"

	log.Printf("Forward reader signup to api %s", url)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		SendJSON(s).
		End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) VerifySignup(token string) (*http.Response, error) {
	url := c.baseURL + "/auth/email/verification/" + token

	log.Printf("Forward reader verification to api %s", url)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) Paywall() (*http.Response, error) {
	url := c.baseURL + "/paywall"

	resp, errs := fetch.New().Get(url).SetBearerAuth(c.key).End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
