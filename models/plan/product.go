package plan

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

type BaseProduct struct {
	ID           string      `db:"product_id"`
	Heading      string      `db:"heading"`
	SmallPrint   null.String `db:"small_print"`
	Tier         enum.Tier   `db:"tier"`
	YearlyPlanID string      `db:"yearly_plan_id"`
}

// RawProduct is the db scan target.
// Description fields needs to be split into arrays by \r\n.
type RawProduct struct {
	BaseProduct
	Description string `db:"description"`
}

// Product describes a product present to user
// on UI.
type Product struct {
	BaseProduct
	Description []string
	Plans       []Plan
}

func NewProduct(raw RawProduct) Product {
	return Product{
		BaseProduct: raw.BaseProduct,
		Description: strings.Split(raw.Description, "\r\n"),
	}
}

// Products is an array of Product.
type Products []Product

func (ps Products) SetPlans(plans Plans) {
	for i, product := range ps {
		if product.Plans == nil {
			ps[i].Plans = make([]Plan, 0)
		}
		plan, ok := plans[product.YearlyPlanID]
		if ok {
			ps[i].Plans = append(ps[i].Plans, plan)
		}
	}
}

func (ps Products) GetPlansIDs() []string {
	var ids = make([]string, 0)

	for _, product := range ps {
		ids = append(ids, product.YearlyPlanID)
	}

	return ids
}
