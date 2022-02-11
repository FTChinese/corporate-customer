package api

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"testing"
)

func TestClient_WxOAuthSession(t *testing.T) {
	faker.MustSetupViper()

	type args struct {
		appID string
	}
	tests := []struct {
		name    string
		fields  Client
		args    args
		want    WxOAuthCodeSession
		wantErr bool
	}{
		{
			name:   "Generate wx oauth session",
			fields: NewSubsAPIClient(false),
			args: args{
				appID: "app123456",
			},
			want:    WxOAuthCodeSession{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			got, err := c.WxOAuthSession(tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WxOAuthCodeSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//    t.Errorf("WxOAuthCodeSession() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
