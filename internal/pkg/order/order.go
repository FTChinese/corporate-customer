package order

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	sq2 "github.com/FTChinese/ftacademy/pkg/sq"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// CartItem is the plan user is subscribing
// and the number of copies for this plan.
//type CartItem struct {
//	PlanID     string             `json:"planId"`
//	Quantity   int64              `json:"quantity"`
//	CycleCount int64              `json:"cycleCount"`
//	Plan       plan2.DiscountPlan `json:"-"`
//}

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

//func (c Cart) BuildOrders(teamID string) OrderList {
//	var orders []Order
//
//	for _, v := range c.Items {
//		for i := 0; i < int(v.Quantity); i++ {
//			o := NewOrder(v, teamID, c.CheckoutID)
//			orders = append(orders, o)
//		}
//	}
//
//	return orders
//}

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
	ID             string          `db:"order_id"`
	AmountPayable  float64         `json:"amountPayable"`
	CreatedBy      string          `json:"createdBy"`
	CreatedUTC     chrono.Time     `json:"createdUtc"`
	Status         Status          `json:"status"`
	ItemCount      int64           `json:"itemCount"`
	CartItemsBrief []CartItemBrief `json:"cart_items_brief"` // An array of products, together with the quantities, use is trying to purchase.
	TeamID         string          `db:"team_id"`
	// The following fields exists only after payment confirmed.
	AmountPaid    null.Float    `json:"amountPaid"`
	ApprovedBy    null.String   `json:"approvedBy"`
	ApprovedUTC   chrono.Time   `json:"approvedUtc"`
	SellerNote    null.String   `json:"sellerNote"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
	TransactionID null.String   `json:"transactionId"`
}

func NewOrder(cart ShoppingCart, p admin.PassportClaims) Order {

	return Order{
		ID:             pkg.OrderID(),
		AmountPayable:  cart.TotalAmount,
		CreatedBy:      p.AdminID,
		CreatedUTC:     chrono.TimeNow(),
		Status:         StatusPending,
		ItemCount:      cart.ItemCount,
		CartItemsBrief: cart.ItemsBrief(),
		TeamID:         p.TeamID.String,
		AmountPaid:     null.Float{},
		ApprovedBy:     null.String{},
		ApprovedUTC:    chrono.Time{},
		SellerNote:     null.String{},
		PaymentMethod:  "",
		TransactionID:  null.String{},
	}
}

//func NewOrder(item CartItem, teamID, checkoutID string) Order {
//	return Order{
//		ID:           "ord_" + rand.String(12),
//		PlanID:       item.PlanID,
//		DiscountID:   item.Plan.DiscountID,
//		LicenceID:    "lic_" + rand.String(12),
//		TeamID:       teamID,
//		CheckoutID:   checkoutID,
//		Amount:       item.Plan.PayableAmount(),
//		CycleCount:   item.CycleCount,
//		TrialDays:    7,
//		Kind:         enum.OrderKindCreate,
//		PeriodStart:  chrono.Date{},
//		PeriodEnd:    chrono.Date{},
//		CreatedUTC:   chrono.TimeNow(),
//		ConfirmedUTC: chrono.Time{},
//	}
//}

func (o Order) RowValues() []interface{} {
	return []interface{}{
		//o.ID,
		//o.PlanID,
		//o.DiscountID,
		//o.LicenceID,
		//o.TeamID,
		//o.CheckoutID,
		//o.Amount,
		//o.CycleCount,
		//o.TrialDays,
		//o.Kind,
		//"UTC_TIMESTAMP()",
	}
}

type OrderList []Order

func (ol OrderList) Each(handler func(row sq2.InsertRow)) {
	for _, o := range ol {
		handler(o)
	}
}

// PagedOrders contains the count of total orders of a team,
// and the current page of orders.
type PagedOrders struct {
	Total int64   `json:"total"`
	Data  []Order `json:"data"`
	Err   error   `json:"-"`
}
