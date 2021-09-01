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
				order: adm.OrderByItems(
					checkout.OrderItem{
						Price:         price.MockPriceStdYear,
						NewCopies:     5,
						RenewalCopies: 3,
					},
					checkout.OrderItem{
						Price:         price.MockPricePrm,
						NewCopies:     2,
						RenewalCopies: 0,
					},
				),
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
						licence.NewLicBuilder(price.MockPriceStdYear).Build(),
						licence.NewLicBuilder(price.MockPriceStdYear).Build(),
					},
				}.Schema(pkg.OrderID()),
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

func TestTxRepo_SaveLicenceQueue(t *testing.T) {
	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		q checkout.BulkLicenceQueue
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
				q: checkout.NewCartBuilder().
					AddStandardN(5).
					AddRenewal(licence.NewLicBuilder(price.MockPriceStdYear).Build()).
					AddPremiumN(2).
					Build().
					LicenceQueue(pkg.OrderID()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.SaveLicenceQueue(tt.args.q); (err != nil) != tt.wantErr {
				t.Errorf("SaveLicenceQueue() error = %v, wantErr %v", err, tt.wantErr)
				_ = tx.Rollback()
				return
			}

			_ = tx.Commit()
		})
	}
}
