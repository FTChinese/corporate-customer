package letter

import (
	"github.com/brianvoe/gofakeit/v5"
	"testing"
)

func TestCtxVerification_Render(t *testing.T) {
	type fields struct {
		Email    string
		UserName string
		Link     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Render verification email",
			fields: fields{
				Email:    gofakeit.Email(),
				UserName: gofakeit.Username(),
				Link:     gofakeit.URL(),
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CtxVerification{
				Email:    tt.fields.Email,
				UserName: tt.fields.UserName,
				Link:     tt.fields.Link,
			}
			got, err := ctx.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("Render() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", got)
		})
	}
}
