package plan

// Discount is the amount to subtract from Plan.Price when
// user purchases licence in bulk.
// For example, when purchase 10 copies, the price of each
// licence will off 10.
// Example:
type Discount struct {
	Quantity int64   `json:"quantity" db:"quantity"`  // The amount of minimum copies of licences purchased when this discount becomes available.
	PriceOff float64 `json:"priceOff" db:"price_off"` // Deducted from Plan.Price
}

// IsZero tests if a the discount actually exist since
// we use a JOIN to retrieve the plan together with its
// discount in one go.
func (d Discount) IsZero() bool {
	return d.PriceOff <= 0
}
