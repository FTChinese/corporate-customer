package subs

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	api2 "github.com/FTChinese/ftacademy/internal/repository/api"
	txrepo2 "github.com/FTChinese/ftacademy/internal/repository/txrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	gorest "github.com/FTChinese/go-rest"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_LoadLicence(t *testing.T) {

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))
	lic := licence.MockLicence(price.MockPriceStdYear)
	txrepo2.MockNewRepo().MustCreateLicence(lic.BaseLicence)

	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Licence
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
	lic := licence.MockLicence(price.MockPriceStdYear)
	txrepo2.MockNewRepo().MustCreateLicence(lic.BaseLicence)

	type args struct {
		teamID string
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []licence.Licence
		wantErr bool
	}{
		{
			name: "list licences",
			args: args{
				teamID: lic.TeamID,
				page:   gorest.NewPagination(1, 10),
			},
			want:    []licence.Licence{lic},
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
	lic := licence.MockLicence(price.MockPriceStdYear)
	txrepo2.MockNewRepo().MustCreateLicence(lic.BaseLicence)

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

func TestEnv_GrantLicence(t *testing.T) {

	a := api2.MockNewClient().MustCreateAssignee()

	mockRepo := txrepo2.MockNewRepo()
	mockRepo.MustCreateMember(reader.MockMembership(a.FtcID.String))
	lic := mockRepo.MustCreateInvitedLicence(a)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		r  admin.AccessRight
		to licence.Assignee
	}
	tests := []struct {
		name    string
		args    args
		want    licence.GrantResult
		wantErr bool
	}{
		{
			name: "Grant licence",
			args: args{
				r: admin.AccessRight{
					RowID:  lic.ID,
					TeamID: lic.TeamID,
				},
				to: a,
			},
			want:    licence.GrantResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.GrantLicence(tt.args.r, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrantLicence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GrantLicence() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_RevokeLicence(t *testing.T) {
	a := api2.MockNewClient().MustCreateAssignee()

	mockRepo := txrepo2.MockNewRepo()

	lic := mockRepo.MustCreateGrantedLicence(a)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		args    args
		want    licence.RevokeResult
		wantErr bool
	}{
		{
			name: "Revoke licence",
			args: args{
				r: admin.AccessRight{
					RowID:  lic.ID,
					TeamID: lic.TeamID,
				},
			},
			want:    licence.RevokeResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RevokeLicence(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RevokeLicence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RevokeLicence() got = %v, want %v", got, tt.want)
			//}
			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
