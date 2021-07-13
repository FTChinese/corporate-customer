package price

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/guregu/null"
)

type Source string

const (
	SourceFTC    = "ftc"
	SourceStripe = "stripe"
)

// Price presents the price of a price. It unified prices coming
// from various source, e.g., FTC in-house or Stripe API.
type Price struct {
	ID string `json:"id"`
	Edition
	Active     bool        `json:"active"`
	Currency   Currency    `json:"currency"`
	LiveMode   bool        `json:"liveMode"`
	Nickname   null.String `json:"nickname"`
	ProductID  string      `json:"productId"`
	Source     Source      `json:"source"`
	UnitAmount float64     `json:"unitAmount"`
}

// Value turns a Price instance to JSON upon saving
func (p Price) Value() (driver.Value, error) {
	if p.ID == "" {
		return nil, nil
	}

	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (p *Price) Scan(src interface{}) error {
	if src == nil {
		*p = Price{}
	}

	switch s := src.(type) {
	case []byte:
		var tmp Price
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*p = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}
