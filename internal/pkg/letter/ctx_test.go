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

func TestCtxPwReset_Render(t *testing.T) {
	type fields struct {
		UserName string
		Link     string
		Duration string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Password reset letter",
			fields: fields{
				UserName: gofakeit.Username(),
				Link:     gofakeit.URL(),
				Duration: "3小时",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CtxPwReset{
				UserName: tt.fields.UserName,
				Link:     tt.fields.Link,
				Duration: tt.fields.Duration,
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

func TestCtxInvitation_Render(t *testing.T) {
	type fields struct {
		ToName     string
		AdminEmail string
		TeamName   string
		Tier       string
		URL        string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Invitation",
			fields: fields{
				ToName:     gofakeit.Username(),
				AdminEmail: gofakeit.Email(),
				TeamName:   gofakeit.Company(),
				Tier:       "Standard",
				URL:        gofakeit.URL(),
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CtxInvitation{
				ToName:     tt.fields.ToName,
				AdminEmail: tt.fields.AdminEmail,
				TeamName:   tt.fields.TeamName,
				Tier:       tt.fields.Tier,
				URL:        tt.fields.URL,
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

func TestCtxLicenceGranted_Render(t *testing.T) {
	type fields struct {
		Name           string
		AssigneeEmail  string
		Tier           string
		ExpirationDate string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Licence granted",
			fields: fields{
				Name:           gofakeit.Username(),
				AssigneeEmail:  gofakeit.Email(),
				Tier:           "Standard",
				ExpirationDate: "2022-12-12",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CtxLicenceGranted{
				Name:           tt.fields.Name,
				AssigneeEmail:  tt.fields.AssigneeEmail,
				Tier:           tt.fields.Tier,
				ExpirationDate: tt.fields.ExpirationDate,
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
