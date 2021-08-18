package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/go-rest/chrono"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
	"time"
)

func TestEnv_CreateTeam(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		t admin.Team
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create team",
			args: args{
				t: admin.MockTeam(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.CreateTeam(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("CreateTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadTeam(t *testing.T) {
	faker.SeedGoFake()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	team := admin.MockTeam()
	team.CreatedUTC = chrono.TimeFrom(team.CreatedUTC.Truncate(time.Second).In(time.UTC))
	_ = env.CreateTeam(team)

	type args struct {
		teamID  string
		adminID string
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
				teamID:  team.ID,
				adminID: team.AdminID,
			},
			want:    team,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadTeam(tt.args.teamID, tt.args.adminID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadTeam() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_UpdateTeam(t *testing.T) {

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	team := admin.MockTeam()

	_ = env.CreateTeam(team)

	type args struct {
		t admin.Team
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update team",
			args: args{
				t: team.Update(admin.MockTeamParams()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateTeam(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
