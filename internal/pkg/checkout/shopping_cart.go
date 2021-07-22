package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
)

// ShoppingCart is used to hold data submitted by client.
// When storing a shopping cart, its data is decomposed:
// A summary of items is save in an order row;
// the items array is save into a separate table one by one.
// Table `order` to `order_item` has a one-many relation.
type ShoppingCart struct {
	Items       []CartItem `json:"items"`
	ItemCount   int64      `json:"itemCount"`
	TotalAmount float64    `json:"totalAmount"`
}

// CartItem describes how many copies user purchased a product.
// Saved in b2b.order_item table with Price and Renewals fields as JSON.
type CartItem struct {
	Price     price.Price     `json:"price" db:"price_snapshot"`
	NewCopies int64           `json:"newCopies" db:"new_copy_count"`
	Renewals  RenewalLicences `json:"renewals" db:"renewal_licences"`
}

// RenewalLicences to ease retrieve/save an array of
// licences into db.
type RenewalLicences []licence.Licence

func (rl RenewalLicences) Value() (driver.Value, error) {
	j, err := json.Marshal(rl)
	if err != nil {
		return nil, err
	}

	return string(j), nil
}

func (rl *RenewalLicences) Scan(src interface{}) error {
	if src == nil {
		*rl = []licence.Licence{}
		return nil
	}
	switch s := src.(type) {
	case []byte:
		var tmp []licence.Licence
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*rl = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to []Licence")
	}
}
