package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/pkg/price"
	"reflect"
	"testing"
)

func TestCartBuilder_AddNewN(t *testing.T) {
	adm := NewAdmin()

	type fields struct {
		admin Admin
		store map[string]checkout.CartItem
	}
	type args struct {
		p price.Price
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CartBuilder
	}{
		{
			name: "Add n new items",
			fields: fields{
				admin: adm,
				store: map[string]checkout.CartItem{},
			},
			args: args{
				p: price.MockPriceStdYear,
				n: 5,
			},
			want: CartBuilder{
				admin: adm,
				store: map[string]checkout.CartItem{
					price.MockPriceStdYear.ID: {
						Price:     price.MockPriceStdYear,
						NewCopies: 5,
						Renewals:  nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := CartBuilder{
				admin: tt.fields.admin,
				store: tt.fields.store,
			}
			if got := b.AddNewN(tt.args.p, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddNewN() = %v, want %v", got, tt.want)
			}
		})
	}
}
