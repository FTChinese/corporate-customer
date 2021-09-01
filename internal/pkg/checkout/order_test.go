package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"testing"
)

func TestNewOrder(t *testing.T) {
	type args struct {
		cart ShoppingCart
		p    admin.PassportClaims
	}
	tests := []struct {
		name string
		args args
		want Order
	}{
		{
			name: "Create order",
			args: args{
				cart: ShoppingCart{
					Items: []CartItem{
						{
							Price:     price.MockPriceStdYear,
							NewCopies: 5,
							Renewals:  nil,
						},
						{
							Price:     price.MockPricePrm,
							NewCopies: 2,
							Renewals:  nil,
						},
					},
					ItemCount:   0,
					TotalAmount: 0,
				},
				p: admin.PassportClaims{
					AdminID: uuid.New().String(),
					TeamID:  null.StringFrom(pkg.TeamID()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOrder(tt.args.cart, tt.args.p)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewOrder() = %v, want %v", got, tt.want)
			//	return
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
