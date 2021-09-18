package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestTxRepo_LockLicenceTxn(t *testing.T) {

	orderSchema := mock.NewAdmin().
		CartBuilder().
		AddNewStandardN(1).
		BuildOrderSchema()

	mock.NewRepo().InsertOrderSchema(orderSchema)

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    checkout.LicenceTransaction
		wantErr bool
	}{
		{
			name: "Lock licence txn",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				id: orderSchema.Transactions[0].ID,
			},
			want:    orderSchema.Transactions[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			got, err := tx.LockLicenceTxn(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LockLicenceTxn() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}
			_ = tx.Commit()

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LockLicenceTxn() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestTxRepo_FinalizeLicenceTxn(t *testing.T) {

	orderSchema := mock.NewAdmin().
		CartBuilder().
		AddNewStandardN(1).
		BuildOrderSchema()

	mock.NewRepo().InsertOrderSchema(orderSchema)

	txn1 := orderSchema.Transactions[0]

	txn1 = txn1.Finalize()

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		lt checkout.LicenceTransaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Licence transaction finalized",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				lt: txn1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.FinalizeLicenceTxn(tt.args.lt); (err != nil) != tt.wantErr {
				t.Errorf("FinalizeLicenceTxn() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_LockLicenceCMS(t *testing.T) {
	lic := mock.NewAdmin().StdLicenceBuilder().Build()

	mock.NewRepo().InsertLicence(lic)

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    licence.Licence
		wantErr bool
	}{
		{
			name: "Lock a licence from cms",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				id: lic.ID,
			},
			//want:    licence.Licence{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			got, err := tx.LockLicenceCMS(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LockLicenceCMS() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LockLicenceCMS() got = %v, want %v", got, tt.want)
			//}

			_ = tx.Commit()

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestTxRepo_CreateLicence(t *testing.T) {

	lic := mock.NewAdmin().StdLicenceBuilder().Build()

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
			name: "Create licence",
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
			if err := tx.CreateLicence(tt.args.lic); (err != nil) != tt.wantErr {
				t.Errorf("CreateLicence() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_RenewLicence(t *testing.T) {
	lic := mock.NewAdmin().StdLicenceBuilder().Build()

	mock.NewRepo().InsertLicence(lic)

	lic = lic.Renewed(price.MockPriceStdYear, pkg.OrderID())

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
			name: "Renew licence",
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
			if err := tx.RenewLicence(tt.args.lic); (err != nil) != tt.wantErr {
				t.Errorf("RenewLicence() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}
