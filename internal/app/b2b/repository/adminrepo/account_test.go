package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/brianvoe/gofakeit/v5"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_BaseAccountByID(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := admin.MockAccount()
	_ = env.SignUp(account)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    admin.BaseAccount
		wantErr bool
	}{
		{
			name: "Retrieve base account",
			args: args{
				id: account.ID,
			},
			want:    account.BaseAccount,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.BaseAccountByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseAccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseAccountByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_BaseAccountByEmail(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := admin.MockAccount()
	_ = env.SignUp(account)

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    admin.BaseAccount
		wantErr bool
	}{
		{
			name: "Base account by email",
			args: args{
				email: account.Email,
			},
			want:    account.BaseAccount,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.BaseAccountByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseAccountByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseAccountByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_UpdateName(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := admin.MockAccount()
	_ = env.SignUp(account)

	faker.SeedGoFake()

	type args struct {
		account admin.BaseAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update name",
			args: args{
				account: account.UpdateName(gofakeit.Username()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateName(tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("UpdateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_UpdatePassword(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := admin.MockAccount()
	_ = env.SignUp(account)

	type args struct {
		p input.PasswordUpdateParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update password",
			args: args{
				p: input.PasswordUpdateParams{
					ID:  account.ID,
					Old: "",
					New: faker.SimplePassword(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.UpdatePassword(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
