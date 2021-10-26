package api

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"testing"
)

func TestClient_LoadAccountByFtcID(t *testing.T) {
	faker.MustSetupViper()

	c := NewSubsAPIClient(false)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    fetch.Response
		wantErr bool
	}{
		{
			name: "Load ftc account",
			args: args{
				id: "405d536a-f926-420a-862a-c9eea060223c",
			},
			want:    fetch.Response{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.LoadAccountByFtcID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAccountByFtcID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadAccountByFtcID() got = %v, want %v", got, tt.want)
			//}

			t.Logf("Status code %d", got.StatusCode)
			t.Logf("Body %s", got.Body)
		})
	}
}
