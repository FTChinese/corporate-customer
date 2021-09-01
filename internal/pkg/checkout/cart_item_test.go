package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"reflect"
	"testing"
)

func TestCartItem_Summary(t *testing.T) {
	type fields struct {
		Price     price.Price
		NewCopies int64
		Renewals  ExpLicenceListJSON
	}
	tests := []struct {
		name   string
		fields fields
		want   OrderItem
	}{
		{
			name: "Convert to order item",
			fields: fields{
				Price:     price.MockPriceStdYear,
				NewCopies: 5,
				Renewals:  nil,
			},
			want: OrderItem{
				Price:         price.MockPriceStdYear,
				NewCopies:     5,
				RenewalCopies: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := CartItem{
				Price:     tt.fields.Price,
				NewCopies: tt.fields.NewCopies,
				Renewals:  tt.fields.Renewals,
			}
			got := ci.OrderItem()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderItem() = %v, want %v", got, tt.want)
				return
			}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
