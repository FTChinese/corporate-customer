package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/price"
	"reflect"
	"testing"
)

func TestShoppingCart_OrderItemList(t *testing.T) {
	type fields struct {
		Items       []CartItem
		ItemCount   int64
		TotalAmount float64
	}
	tests := []struct {
		name   string
		fields fields
		want   OrderItemListJSON
	}{
		{
			name: "convert cart items to order items",
			fields: fields{
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
				ItemCount:   7,
				TotalAmount: 5*price.MockPriceStdYear.UnitAmount + 2*price.MockPricePrm.UnitAmount,
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
			c := ShoppingCart{
				Items:       tt.fields.Items,
				ItemCount:   tt.fields.ItemCount,
				TotalAmount: tt.fields.TotalAmount,
			}
			if got := c.OrderItemList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderItemList() = %v, want %v", got, tt.want)
			}
		})
	}
}
