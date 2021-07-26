package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
)

func TestTxRepo_RetrieveBaseLicence(t *testing.T) {
	lic := licence.MockLicence(price.MockPriceStdYear)

	err := MockNewRepo().MockCreateLicence(lic)
	if err != nil {
		t.Error(err)
		return
	}

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
		want    licence.BaseLicence
		wantErr bool
	}{
		{
			name: "Retrieve a licence",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				r: admin.AccessRight{
					RowID:  lic.ID,
					TeamID: lic.TeamID,
				},
			},
			want:    lic.BaseLicence,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			got, err := tx.RetrieveBaseLicence(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveBaseLicence() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveBaseLicence() got = %v, \nwant %v", got, tt.want)
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_UpdateLicenceStatus(t *testing.T) {
	lic := licence.MockLicence(price.MockPriceStdYear)

	err := MockNewRepo().MockCreateLicence(lic)
	if err != nil {
		t.Error(err)
		return
	}

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		lic licence.BaseLicence
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Licence with invitation",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				lic: lic.WithInvitation(licence.MockInvitation(lic)),
			},
			wantErr: false,
		},
		{
			name: "Licence invitation revoked",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				lic: lic.WithInvitationRevoked(),
			},
			wantErr: false,
		},
		{
			name: "Licence granted",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				lic: lic.Granted(licence.MockAssignee(), licence.MockInvitation(lic)),
			},
			wantErr: false,
		},
		{
			name: "Licence revoked",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				lic: lic.Revoked(),
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
				return
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_CreateInvitation(t *testing.T) {
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
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				inv: licence.MockInvitation(licence.MockLicence(price.MockPriceStdYear)),
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
