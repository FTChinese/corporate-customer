package api

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"net/http"
	"testing"
)

func TestClient_Paywall(t *testing.T) {
	faker.MustSetupViper()

	tests := []struct {
		name    string
		client  Client
		want    *http.Response
		wantErr bool
	}{
		{
			name:    "Production paywall",
			client:  NewSubsAPIClient(true),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.client.Paywall()
			if (err != nil) != tt.wantErr {
				t.Errorf("Paywall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Stats=%d", got.StatusCode)
			t.Logf("Body=%s", faker.MustReadBody(got.Body))

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Paywall() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
