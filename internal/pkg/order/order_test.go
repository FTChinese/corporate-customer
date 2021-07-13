package order

import (
	sq2 "github.com/FTChinese/ftacademy/pkg/sq"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"testing"
)

func TestNewOrder(t *testing.T) {
	var orders OrderList = make([]Order, 0)
	for i := 0; i < 10; i++ {
		o := Order{
			ID:           "ord_" + rand.String(12),
			PlanID:       "plan_" + rand.String(12),
			DiscountID:   null.Int{},
			LicenceID:    "lic_" + rand.String(12),
			TeamID:       "team_" + rand.String(12),
			CheckoutID:   "chk_" + rand.String(12),
			Amount:       258,
			CycleCount:   1,
			TrialDays:    7,
			Kind:         enum.OrderKindCreate,
			PeriodStart:  chrono.Date{},
			PeriodEnd:    chrono.Date{},
			CreatedUTC:   chrono.TimeNow(),
			ConfirmedUTC: chrono.Time{},
		}

		orders = append(orders, o)
	}

	values := sq2.BuildInsertValues(orders)

	t.Logf("%v", values)
}
