package subs

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	api2 "github.com/FTChinese/ftacademy/internal/repository/api"
	txrepo2 "github.com/FTChinese/ftacademy/internal/repository/txrepo"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_RetrieveAssignee(t *testing.T) {
	faker.SeedGoFake()

	a := api2.MockNewClient().MustCreateAssignee()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Assignee
		wantErr bool
	}{
		{
			name: "Retrieve assignee",
			args: args{
				id: a.FtcID.String,
			},
			want:    a,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.RetrieveAssignee(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAssignee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveAssignee() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_FindAssignee(t *testing.T) {
	faker.SeedGoFake()

	a := api2.MockNewClient().MustCreateAssignee()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    licence.Assignee
		wantErr bool
	}{
		{
			name: "Find assignee",
			args: args{
				email: a.Email.String,
			},
			want:    a,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.FindAssignee(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAssignee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAssignee() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_RetrieveMembership(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	m := reader.MockMembership("")

	txrepo2.MockNewRepo().MustCreateMember(m)

	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		args    args
		want    reader.Membership
		wantErr bool
	}{
		{
			name: "Retrieve membership",
			args: args{
				compoundID: m.CompoundID,
			},
			want:    m,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveMembership(tt.args.compoundID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveMembership() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_ArchiveMembership(t *testing.T) {
	m := reader.MockMembership("")
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		m reader.MemberSnapshot
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Archive membership",
			args: args{
				m: m.Archive(reader.B2BArchiver(reader.ArchiveActionGrant)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.ArchiveMembership(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ArchiveMembership() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
