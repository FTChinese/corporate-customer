package reader

import (
	"net/url"
	"testing"
)

func TestWxOAuthCodeParams_EncodeQuery(t *testing.T) {
	type fields struct {
		AppID        string
		RedirectURI  string
		ResponseType string
		Scope        string
		State        string
		Fragment     string
	}
	tests := []struct {
		name    string
		fields  wxOAuthCodeParams
		want    string
		wantErr bool
	}{
		{
			name:    "Encode query",
			fields:  newWxOAuthCodeParams("anything"),
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields
			got, err := p.encodeQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("encodeQuery() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", got)
		})
	}
}

func TestWxOAuthCodeParams_URL(t *testing.T) {
	tests := []struct {
		name    string
		fields  wxOAuthCodeParams
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "Oauth url",
			fields:  newWxOAuthCodeParams("whatever"),
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields
			got, err := p.url()
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
