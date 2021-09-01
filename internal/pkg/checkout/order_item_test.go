package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"reflect"
	"testing"
)

func Test_newOrderItemList(t *testing.T) {
	type args struct {
		items []CartItem
	}
	tests := []struct {
		name string
		args args
		want OrderItemListJSON
	}{
		{
			name: "convert cart items to order items",
			args: args{
				items: []CartItem{
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
			},
			want: OrderItemListJSON{
				{
					Price:         price.MockPriceStdYear,
					NewCopies:     5,
					RenewalCopies: 0,
				},
				{
					Price:         price.MockPricePrm,
					NewCopies:     2,
					RenewalCopies: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newOrderItemList(tt.args.items)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOrderItemList() = %v, want %v", got, tt.want)
				return
			}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})

	}
}
