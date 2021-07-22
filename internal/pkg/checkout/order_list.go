package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
)

// OrderedProduct describes the quantity of a product put into an order.
// This is used when save all items of an order as JSON
// in an order's row.
type OrderedProduct struct {
	Price         price.Price `json:"price"`
	NewCopies     int64       `json:"newCopies"`     // How many new copies user purchased
	RenewalCopies int64       `json:"renewalCopies"` // How many renewals user purchased.
}

// OrderedProducts is used the retrieve/save an array of
// CheckoutProduct into db.
type OrderedProducts []OrderedProduct

// Value implements Valuer interface when saving
func (b OrderedProducts) Value() (driver.Value, error) {
	j, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return string(j), nil
}

func (b *OrderedProducts) Scan(src interface{}) error {
	if src == nil {
		*b = []OrderedProduct{}
		return nil
	}
	switch s := src.(type) {
	case []byte:
		var tmp []OrderedProduct
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*b = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to []CheckoutProduct")
	}
}

// BriefOrder describes the details of each transaction
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
type BriefOrder struct {
	BaseOrder
	// An array of products, together with the quantities, use is trying to purchase.
	Products OrderedProducts `json:"products" db:"checkout_products"`
}

func NewBriefOrder(cart ShoppingCart, p admin.PassportClaims) BriefOrder {
	return BriefOrder{
		BaseOrder: BaseOrder{
			ID:            pkg.OrderID(),
			AmountPayable: cart.TotalAmount,
			CreatedBy:     p.AdminID,
			CreatedUTC:    chrono.TimeNow(),
			ItemCount:     cart.ItemCount,
			Status:        StatusPending,
			TeamID:        p.TeamID.String,
		},
		Products: cart.ProductsBrief(),
	}
}

// BriefOrderList contains a list of orders
type BriefOrderList struct {
	pkg.PagedList
	Data []BriefOrder `json:"data"`
}
