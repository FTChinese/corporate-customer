package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/price"
	"reflect"
	"testing"
)

func TestCartBuilder_Add(t *testing.T) {
	type fields struct {
		store map[string]CartItem
	}
	type args struct {
		p price.Price
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CartBuilder
	}{
		{
			name: "Add a new subscription",
			fields: fields{
				store: map[string]CartItem{},
			},
			args: args{
				p: price.MockPriceStdYear,
			},
			want: CartBuilder{
				store: map[string]CartItem{
					price.MockPriceStdYear.ID: {
						Price:     price.MockPriceStdYear,
						NewCopies: 1,
						Renewals:  nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := CartBuilder{
				store: tt.fields.store,
			}
			if got := b.Add(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
