package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/adminrepo"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_AdminProfile(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	account := admin.MockAccount()

	_ = adminrepo.NewEnv(db.MockMySQL(), zaptest.NewLogger(t)).
		SignUp(account)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    admin.Profile
		wantErr bool
	}{
		{
			name: "Admin profile",
			args: args{
				id: account.ID,
			},
			want: admin.Profile{
				BaseAccount: account.BaseAccount,
				TeamParams:  input.TeamParams{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.AdminProfile(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminProfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
