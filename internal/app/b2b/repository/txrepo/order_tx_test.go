package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestTxRepo_CreateOrder(t *testing.T) {
	myDB := db.MockMySQL().Write

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		orderRow checkout.OrderRow
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
				Tx: myDB.MustBegin(),
			},
			args: args{
				orderRow: checkout.MockOrderInputSchema().OrderRow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			if err := tx.CreateOrder(tt.args.orderRow); (err != nil) != tt.wantErr {
				_ = tx.Rollback()
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			_ = tx.Commit()
		})
	}
}

func TestTxRepo_CreateOrderItem(t *testing.T) {
	myDB := db.MockMySQL().Write

	type fields struct {
		Tx *sqlx.Tx
	}
	type args struct {
		item checkout.OrderItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create order item",
			fields: fields{
				Tx: myDB.MustBegin(),
			},
			args: args{
				item: checkout.MockOrderInputSchema().ItemRows[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxRepo{
				Tx: tt.fields.Tx,
			}
			t.Logf("Inserting value %s", faker.MustMarshalIndent(tt.args.item))
			if err := tx.CreateOrderItem(tt.args.item); (err != nil) != tt.wantErr {
				_ = tx.Rollback()
				t.Errorf("CreateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			_ = tx.Commit()
		})
	}
}
