package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/txrepo"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	gorest "github.com/FTChinese/go-rest"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_InvitationByToken(t *testing.T) {

	inv := licence.MockInvitation(licence.MockLicence(price.MockPriceStdYear))

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	txrepo.MockNewRepo().MustCreateInvitation(inv)

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Invitation
		wantErr bool
	}{
		{
			name: "Invitation by token",
			args: args{
				token: inv.Token,
			},
			want:    inv,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.InvitationByToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvitationByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvitationByToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_InvitationByID(t *testing.T) {
	inv := licence.MockInvitation(licence.MockLicence(price.MockPriceStdYear))

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	txrepo.MockNewRepo().MustCreateInvitation(inv)

	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Invitation
		wantErr bool
	}{
		{
			name: "Invitation by id",
			args: args{
				r: admin.AccessRight{
					RowID:  inv.ID,
					TeamID: inv.TeamID,
				},
			},
			want:    inv,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.InvitationByID(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvitationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvitationByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_CreateInvitation(t *testing.T) {
	faker.SeedGoFake()

	lic := licence.MockLicence(price.MockPriceStdYear)

	txrepo.MockNewRepo().MustCreateLicence(lic.Licence)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		params input.InvitationParams
		p      admin.PassportClaims
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Licence
		wantErr bool
	}{
		{
			name: "Create invitation",
			args: args{
				params: input.InvitationParams{
					Email:       gofakeit.Email(),
					Description: null.String{},
					LicenceID:   lic.ID,
				},
				p: admin.PassportClaims{
					AdminID: lic.AdminID,
					TeamID:  null.StringFrom(lic.TeamID),
				},
			},
			want:    licence.Licence{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateInvitation(tt.args.params, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateInvitation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("CreateInvitation() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_RevokeInvitation(t *testing.T) {
	faker.SeedGoFake()

	lic := licence.MockLicence(price.MockPriceStdYear)
	inv := licence.MockInvitation(lic)

	mockRepo := txrepo.MockNewRepo()
	mockRepo.MustCreateLicence(lic.WithInvitation(inv))
	mockRepo.MustCreateInvitation(inv)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		invID  string
		teamID string
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Invitation
		wantErr bool
	}{
		{
			name: "Revoke invitation",
			args: args{
				invID:  inv.ID,
				teamID: inv.TeamID,
			},
			want:    licence.Invitation{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RevokeInvitation(tt.args.invID, tt.args.teamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RevokeInvitation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RevokeInvitation() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_countInvitation(t *testing.T) {
	inv := licence.MockInvitation(licence.MockLicence(price.MockPriceStdYear))
	txrepo.MockNewRepo().MustCreateInvitation(inv)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

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
			name: "count invitation",
			args: args{
				teamID: inv.TeamID,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countInvitation(tt.args.teamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("countInvitation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countInvitation() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_listInvitations(t *testing.T) {
	inv := licence.MockInvitation(licence.MockLicence(price.MockPriceStdYear))
	txrepo.MockNewRepo().MustCreateInvitation(inv)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		teamID string
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []licence.Invitation
		wantErr bool
	}{
		{
			name: "list invitation",
			args: args{
				teamID: inv.TeamID,
				page:   gorest.NewPagination(1, 10),
			},
			want:    []licence.Invitation{inv},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.listInvitations(tt.args.teamID, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("listInvitations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listInvitations() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_ListInvitations(t *testing.T) {
	inv := licence.MockInvitation(licence.MockLicence(price.MockPriceStdYear))
	txrepo.MockNewRepo().MustCreateInvitation(inv)

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		teamID string
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    licence.InvitationList
		wantErr bool
	}{
		{
			name: "List invitations",
			args: args{
				teamID: inv.TeamID,
				page:   gorest.NewPagination(1, 10),
			},
			want:    licence.InvitationList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListInvitations(tt.args.teamID, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListInvitations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ListInvitations() got = %v, want %v", got, tt.want)
			//}
			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
