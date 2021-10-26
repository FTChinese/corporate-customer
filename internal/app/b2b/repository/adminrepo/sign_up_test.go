package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/db"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_SignUp(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		a admin.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sign up",
			args: args{
				a: mock.NewAdmin().Account,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SignUp(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("EmailSignUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
