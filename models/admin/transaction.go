package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Transaction describes the details of each transaction
// to purchase a licence.
// Only when a Transaction is confirmed shall we
// create/renew/upgrade a licence.
type Transaction struct {
	ID           string         `db:"trx_id"`
	Price        float64        `db:"price"`
	Amount       float64        `db:"amount"`
	Tier         enum.Tier      `db:"tier"`
	Cycle        enum.Cycle     `db:"cycle"`
	CycleCount   int64          `db:"cycle_count"`
	TrialDays    int64          `db:"trial_days"`
	PeriodStart  chrono.Date    `db:"period_start"`
	PeriodEnd    chrono.Date    `db:"period_end"`
	Kind         enum.OrderKind `db:"kind"`
	CreatedUTC   chrono.Time    `db:"created_utc"`
	ConfirmedUTC chrono.Time    `db:"confirmed_utc"`
	LicenceID    null.String    `db:"licence_id"`
}
