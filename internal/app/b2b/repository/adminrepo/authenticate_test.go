package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_Authenticate(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := mock.NewAdmin().Account
	_ = env.SignUp(account)

	type args struct {
		params input.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    admin.AuthResult
		wantErr bool
	}{
		{
			name: "Authenticate password",
			args: args{
				params: input.Credentials{
					Email:    account.Email,
					Password: account.Password,
				},
			},
			want: admin.AuthResult{
				AdminID:         account.ID,
				PasswordMatched: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.Authenticate(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_VerifyPassword(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := mock.NewAdmin().Account
	_ = env.SignUp(account)

	type args struct {
		params input.PasswordUpdateParams
	}
	tests := []struct {
		name    string
		args    args
		want    admin.AuthResult
		wantErr bool
	}{
		{
			name: "Verify password by id",
			args: args{
				params: input.PasswordUpdateParams{
					ID:  account.ID,
					Old: account.Password,
					New: "",
				},
			},
			want: admin.AuthResult{
				AdminID:         account.ID,
				PasswordMatched: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.VerifyPassword(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
