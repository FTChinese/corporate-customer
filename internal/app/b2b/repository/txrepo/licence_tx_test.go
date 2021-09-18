package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestTxRepo_LockLicence(t *testing.T) {
	adm := mock.NewAdmin()
	lic := adm.StdLicenceBuilder().Build()

	mock.NewRepo().InsertLicence(lic)

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    licence.Licence
		wantErr bool
	}{
		{
			name: "Lock licence",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				r: admin.AccessRight{
					RowID:  lic.ID,
					TeamID: adm.TeamID.String,
				},
			},
			want:    licence.Licence{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			got, err := tx.LockLicence(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("LockLicence() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LockLicence() got = %v, want %v", got, tt.want)
			//}
			_ = tx.Commit()
			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestTxRepo_UpdateLicenceStatus(t *testing.T) {
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		Build()

	mock.NewRepo().InsertLicence(lic.Revoked())

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		lic licence.Licence
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update licence status",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				lic: lic,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.UpdateLicenceStatus(tt.args.lic); (err != nil) != tt.wantErr {
				t.Errorf("UpdateLicenceStatus() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_CreateInvitation(t *testing.T) {
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		Build()

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		inv licence.Invitation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create invitation",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				inv: lic.LatestInvitation.Invitation,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.CreateInvitation(tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("CreateInvitation() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_RetrieveInvitation(t *testing.T) {
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		Build()

	inv := lic.LatestInvitation.Invitation

	mock.NewRepo().InsertInvitation(inv)

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    licence.Invitation
		wantErr bool
	}{
		{
			name: "Retrieve invitation",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				r: admin.AccessRight{
					RowID:  inv.ID,
					TeamID: inv.TeamID,
				},
			},
			want:    licence.Invitation{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			got, err := tx.RetrieveInvitation(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveInvitation() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RetrieveInvitation() got = %v, want %v", got, tt.want)
			//}
			_ = tx.Commit()

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestTxRepo_UpdateInvitationStatus(t *testing.T) {
	lic := mock.NewAdmin().
		StdLicenceBuilder().
		SetPersona(mock.NewPersona()).
		Build()

	inv := lic.LatestInvitation.Invitation

	mock.NewRepo().InsertInvitation(inv)

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		inv licence.Invitation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update invitation status",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				inv: inv.Accepted(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.UpdateInvitationStatus(tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInvitationStatus() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}
