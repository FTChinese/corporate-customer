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
		want    WxOAuthSession
		wantErr bool
	}{
		{
			name:   "Generate wx oauth session",
			fields: NewSubsAPIClient(false),
			args: args{
				appID: "app123456",
			},
			want:    WxOAuthSession{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			got, err := c.WxOAuthSession(tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WxOAuthSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//    t.Errorf("WxOAuthSession() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
