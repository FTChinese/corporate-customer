package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/go-rest/chrono"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
	"time"
)

func TestEnv_SaveEmailVerifier(t *testing.T) {
	faker.SeedGoFake()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		v admin.EmailVerifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save email verification",
			args: args{
				v: mock.NewAdmin().EmailVerifier(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.SaveEmailVerifier(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SaveEmailVerifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_RetrieveEmailVerifier(t *testing.T) {
	faker.SeedGoFake()
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	v := mock.NewAdmin().EmailVerifier()
	// Hack to make time equal.
	v.CreatedUTC = chrono.TimeFrom(v.CreatedUTC.Truncate(time.Second).In(time.UTC))
	_ = env.SaveEmailVerifier(v)

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    admin.EmailVerifier
		wantErr bool
	}{
		{
			name: "Retrieve email verifier",
			args: args{
				token: v.Token,
			},
			want:    v,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveEmailVerifier(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveEmailVerifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveEmailVerifier() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_EmailVerified(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := mock.NewAdmin().Account
	_ = env.SignUp(account)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Email verified",
			args: args{
				id: account.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.EmailVerified(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("EmailVerified() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
