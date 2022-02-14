package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

func TestTxRepo_CreateMember(t *testing.T) {
	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		m reader.Membership
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Save membership",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				m: reader.MockMembership(""),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.CreateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("MustCreateMember() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_LockMember(t *testing.T) {

	m := mock.NewPersona().MemberBuilderFTC().Build()

	mock.NewRepo().InsertMembership(m)

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    reader.Membership
		wantErr bool
	}{
		{
			name: "Load membership",
			fields: fields{
				Tx: db.MockMySQL().Read.MustBegin(),
			},
			args: args{
				compoundID: m.CompoundID,
			},
			want: m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			got, err := tx.LockMember(tt.args.compoundID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LockMember() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LockMember() got = %v, want %v", got, tt.want)
				_ = tx.Rollback()
				return
			}
			_ = tx.Commit()
		})
	}
}

func TestTxRepo_UpdateMember(t *testing.T) {
	m := mock.NewPersona().MemberBuilderFTC().Build()

	mock.NewRepo().InsertMembership(m)

	m.ExpireDate = chrono.DateUTCFrom(time.Now().AddDate(0, 0, -1))
	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		m reader.Membership
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update member",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				m: m,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.UpdateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMember() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_SaveInvoice(t *testing.T) {
	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		inv reader.Invoice
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Save invoice",
			fields: fields{
				Tx: db.MockMySQL().Write.MustBegin(),
			},
			args: args{
				inv: reader.MockMembership("").
					CarryOverInvoice().
					WithLicTxID(null.StringFrom(ids.TxnID())),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.SaveInvoice(tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("SaveInvoice() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}
