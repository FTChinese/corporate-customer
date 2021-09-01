package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
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
