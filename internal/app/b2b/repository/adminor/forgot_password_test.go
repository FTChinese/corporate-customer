package adminor

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
	"time"
)

func TestEnv_SavePwResetSession(t *testing.T) {
	faker.SeedGoFake()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		s admin.PwResetSession
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save password reset session",
			args: args{
				s: admin.MockPwResetSession(input.ForgotPasswordParams{
					Email: gofakeit.Email(),
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.SavePwResetSession(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SavePwResetSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_PwResetSession(t *testing.T) {
	faker.SeedGoFake()
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	s := admin.MockPwResetSession(input.ForgotPasswordParams{
		Email: gofakeit.Email(),
	})
	s.CreatedUTC = chrono.TimeFrom(s.CreatedUTC.Truncate(time.Second).In(time.UTC))

	_ = env.SavePwResetSession(s)

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    admin.PwResetSession
		wantErr bool
	}{
		{
			name: "Password reset session",
			args: args{
				token: s.Token,
			},
			want:    s,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.PwResetSession(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("PwResetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PwResetSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_DisablePasswordReset(t *testing.T) {
	faker.SeedGoFake()
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	s := admin.MockPwResetSession(input.ForgotPasswordParams{
		Email: gofakeit.Email(),
	})

	_ = env.SavePwResetSession(s)

	type args struct {
		s admin.PwResetSession
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Disable password reset",
			args: args{
				s: s.WithUsed(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DisablePasswordReset(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("DisablePasswordReset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
