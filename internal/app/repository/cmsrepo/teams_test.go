package cmsrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_LoadTeam(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	team := mock.NewAdmin().Team()

	mock.NewRepo().InsertTeam(team)

	type args struct {
		teamID string
	}
	tests := []struct {
		name    string
		args    args
		want    admin.Team
		wantErr bool
	}{
		{
			name: "Load team",
			args: args{
				teamID: team.ID,
			},
			want:    admin.Team{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadTeam(tt.args.teamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadTeam() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
