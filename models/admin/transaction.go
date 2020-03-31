package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

// Cart represents a user's purchase action.
type Cart struct {
	PlanID   string `form:"planId"`   // Which plan to buy
	Quantity int64  `form:"quantity"` // How many copies of licence for this plan
}

// Transaction describes the details of each transaction
// to purchase a licence.
// If a transaction is used to purchase a new licence, the
// licence should be created together with the order but marked
// as inactive. Once the transaction is confirmed,
// the licence will be activated and the admin is allowed to
// invite someone to use this licence.
// If a transaction is used to renew/upgrade a licence,
// the licence associated with it won't be touched until
// it is confirmed, which will result licence extended or
// upgraded and the membership (if the licence is granted
// to someone) will be backed up and updated corresponding.
type Transaction struct {
	ID           string         `db:"order_id"`
	PlanID       string         `db:"plan_id"`
	LicenceID    string         `db:"licence_id"`
	TeamID       string         `db:"team_id"`
	Amount       float64        `db:"amount"`
	CycleCount   int64          `db:"cycle_count"`
	TrialDays    int64          `db:"trial_days"`
	PeriodStart  chrono.Date    `db:"period_start"`
	PeriodEnd    chrono.Date    `db:"period_end"`
	Kind         enum.OrderKind `db:"kind"`
	CreatedUTC   chrono.Time    `db:"created_utc"`
	ConfirmedUTC chrono.Time    `db:"confirmed_utc"`
}
