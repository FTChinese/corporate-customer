package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestTxRepo_CreateOrder(t *testing.T) {
	adm := mock.NewAdmin()

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		order checkout.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create order",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				order: adm.CartBuilder().
					AddNewStandardN(10).
					AddRenewalStandardN(10).
					AddNewPremiumN(5).
					AddRenewalPremiumN(5).
					BuildOrderSchema().
					OrderRow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}

			t.Logf("%v", tt.args.order)
			t.Logf("%s", faker.MustMarshalIndent(tt.args.order))

			if err := tx.CreateOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_SaveCartItem(t *testing.T) {

	adm := mock.NewAdmin()

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		c checkout.CartItemSchema
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Save car item",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				c: checkout.CartItem{
					Price:     price.MockPriceStdYear,
					NewCopies: 5,
					Renewals: checkout.ExpLicenceListJSON{
						adm.StdLicenceBuilder().BuildExpanded(),
						adm.StdLicenceBuilder().BuildExpanded(),
					},
				}.Schema(ids.OrderID(), mock.NewAdmin().Creator()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.SaveCartItem(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SaveCartItem() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_SaveLicenceTxnList(t *testing.T) {

	adm := mock.NewAdmin()

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		list checkout.BulkLicenceTxn
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Save licence queue",
			fields: fields{
				Tx: db.MockTx(),
			},
			args: args{
				list: checkout.NewOrderSchemaBuilder(
					adm.CartBuilder().
						AddNewStandardN(5).
						AddRenewalStandardN(5).
						AddNewPremiumN(2).
						Build(),
					adm.PassportClaims(),
				).TransactionList(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.SaveLicenceTxnList(tt.args.list); (err != nil) != tt.wantErr {
				t.Errorf("SaveLicenceTxnList() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}
