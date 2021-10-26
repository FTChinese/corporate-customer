package api

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"testing"
)

func TestClient_RequestLoginSMS(t *testing.T) {
	faker.MustSetupViper()

	c := NewSubsAPIClient(false)

	resp, err := c.RequestLoginSMS(faker.MustMarshalToReader(map[string]string{
		"mobile": "15011481214",
	}))

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%s", faker.MustReadBody(resp.Body))
}

func TestClient_VerifyLoginSMS(t *testing.T) {
	faker.MustSetupViper()

	c := NewSubsAPIClient(false)

	resp, err := c.VerifyLoginSMS(faker.MustMarshalToReader(map[string]string{
		"mobile": "15011481214",
		"code":   "210003",
	}))

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%s", resp.Body)
}
