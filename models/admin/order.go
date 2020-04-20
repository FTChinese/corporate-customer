package admin

import (
	"github.com/FTChinese/b2b/models/plan"
	"github.com/FTChinese/b2b/models/sq"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
)

// GetCartPlansIDs collects the plan ids of items
// in a cart and use it to retrieve the plans
// using FIND_IN_SET
func GetCartPlanIDs(items []CartItem) []string {
	var ids = make([]string, 0)

	for _, item := range items {
		ids = append(ids, item.PlanID)
	}

	return ids
}

// CartItem is the plan user is subscribing
// and the number of copies for this plan.
type CartItem struct {
	PlanID     string            `json:"planId"`
	Quantity   int64             `json:"quantity"`
	CycleCount int64             `json:"cycleCount"`
	Plan       plan.DiscountPlan `json:"-"`
}

type CheckoutItem struct {
	CartItem
}

type CheckoutCounter struct {
	ID         string
	Items      []CheckoutItem
	Total      float64
	CreatedUTC chrono.Time `db:"created_utc"`
}

type Cart struct {
	CheckoutID string
	Items      []CartItem
}

func NewCart(items []CartItem, plans plan.GroupedPlans) (Cart, *render.ValidationError) {
	cart := Cart{
		CheckoutID: "chk_" + rand.String(12),
		Items:      nil,
	}

	for _, v := range items {
		p, ok := plans[v.PlanID]
		if !ok {
			return Cart{}, &render.ValidationError{
				Message: "Missing required field",
				Field:   "planId",
				Code:    render.CodeMissing,
			}
		}

		dp := p.DiscountPlan(v.Quantity)

		v.Plan = dp

		cart.Items = append(cart.Items, v)
	}

	return cart, nil
}

func (c Cart) BuildOrders(teamID string) OrderList {
	var orders []Order

	for _, v := range c.Items {
		for i := 0; i < int(v.Quantity); i++ {
			o := NewOrder(v, teamID, c.CheckoutID)
			orders = append(orders, o)
		}
	}

	return orders
}

// Order describes the details of each transaction
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
type Order struct {
	ID           string         `db:"order_id"`
	PlanID       string         `db:"plan_id"`
	DiscountID   null.Int       `db:"discount_id"`
	LicenceID    string         `db:"licence_id"`
	TeamID       string         `db:"team_id"`
	CheckoutID   string         `db:"checkout_id"`
	Amount       float64        `db:"amount"`
	CycleCount   int64          `db:"cycle_count"`
	TrialDays    int64          `db:"trial_days"`
	Kind         enum.OrderKind `db:"kind"`
	PeriodStart  chrono.Date    `db:"period_start"`
	PeriodEnd    chrono.Date    `db:"period_end"`
	CreatedUTC   chrono.Time    `db:"created_utc"`
	ConfirmedUTC chrono.Time    `db:"confirmed_utc"`
}

func NewOrder(item CartItem, teamID, checkoutID string) Order {
	return Order{
		ID:           "ord_" + rand.String(12),
		PlanID:       item.PlanID,
		DiscountID:   item.Plan.DiscountID,
		LicenceID:    "lic_" + rand.String(12),
		TeamID:       teamID,
		CheckoutID:   checkoutID,
		Amount:       item.Plan.PayableAmount(),
		CycleCount:   item.CycleCount,
		TrialDays:    7,
		Kind:         enum.OrderKindCreate,
		PeriodStart:  chrono.Date{},
		PeriodEnd:    chrono.Date{},
		CreatedUTC:   chrono.TimeNow(),
		ConfirmedUTC: chrono.Time{},
	}
}

func (o Order) RowValues() []interface{} {
	return []interface{}{
		o.ID,
		o.PlanID,
		o.DiscountID,
		o.LicenceID,
		o.TeamID,
		o.CheckoutID,
		o.Amount,
		o.CycleCount,
		o.TrialDays,
		o.Kind,
		"UTC_TIMESTAMP()",
	}
}

type OrderList []Order

func (ol OrderList) Each(handler func(row sq.InsertRow)) {
	for _, o := range ol {
		handler(o)
	}
}

// PageOrders contains the count of total orders of a team,
// and the current page of orders.
type PagedOrders struct {
	Total int64   `json:"total"`
	Data  []Order `json:"data"`
	Err   error   `json:"-"`
}
