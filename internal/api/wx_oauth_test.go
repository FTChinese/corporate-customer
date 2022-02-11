package api

import (
	"net/url"
	"testing"
)

func TestWxOAuthCodeParams_EncodeQuery(t *testing.T) {

	tests := []struct {
		name    string
		fields  WxOAuthCodeRequest
		want    string
		wantErr bool
	}{
		{
			name:    "Encode query",
			fields:  NewWxOAuthCodeRequest("anything", "http://callback.example"),
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields
			got, err := p.EncodeQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("EncodeQuery() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", got)
		})
	}
}

func TestWxOAuthCodeParams_Build(t *testing.T) {
	tests := []struct {
		name    string
		fields  WxOAuthCodeRequest
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "Oauth url",
			fields:  NewWxOAuthCodeRequest("whatever", "http://callback.example"),
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields
			got, err := p.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//    t.Errorf("URL() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", got)
		})
	}
}
