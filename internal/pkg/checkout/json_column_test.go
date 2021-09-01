package checkout

import (
	"database/sql/driver"
	"github.com/FTChinese/ftacademy/pkg/price"
	"testing"
)

func TestOrderItemListJSON_Value(t *testing.T) {
	tests := []struct {
		name    string
		l       OrderItemListJSON
		want    driver.Value
		wantErr bool
	}{
		{
			name: "Save order items as json string",
			l: OrderItemListJSON{
				{
					Price:         price.MockPriceStdYear,
					NewCopies:     5,
					RenewalCopies: 4,
				},
				{
					Price:         price.MockPricePrm,
					NewCopies:     3,
					RenewalCopies: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Value() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%v", got)
		})
	}
}
