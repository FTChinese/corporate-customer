package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/db"
	gorest "github.com/FTChinese/go-rest"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_LoadLicence(t *testing.T) {

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		BuildExpanded()

	mock.NewRepo().InsertLicence(lic.Licence)

	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		args    args
		want    licence.ExpandedLicence
		wantErr bool
	}{
		{
			name: "Load licence",
			args: args{
				r: admin.AccessRight{
					RowID:  lic.ID,
					TeamID: lic.TeamID,
				},
			},
			want:    lic,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadLicence(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadLicence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadLicence() got = %v, \nwant %v", got, tt.want)
			}
		})
	}
}

func TestEnv_listLicences(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		BuildExpanded()

	mock.NewRepo().InsertLicence(lic.Licence)

	type args struct {
		teamID string
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []licence.ExpandedLicence
		wantErr bool
	}{
		{
			name: "list licences",
			args: args{
				teamID: lic.TeamID,
				page:   gorest.NewPagination(1, 10),
			},
			want:    []licence.ExpandedLicence{lic},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.listLicences(tt.args.teamID, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("listLicences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listLicences() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_countLicences(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		BuildExpanded()

	mock.NewRepo().InsertLicence(lic.Licence)

	type args struct {
		teamID string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "count licence",
			args: args{
				teamID: lic.TeamID,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countLicences(tt.args.teamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("countLicences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countLicences() got = %v, want %v", got, tt.want)
			}
		})
	}
}
