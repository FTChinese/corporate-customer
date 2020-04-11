package plan

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

type BaseProduct struct {
	ID           string      `db:"product_id"`
	Tier         enum.Tier   `db:"tier"`
	Heading      string      `db:"heading"`
	SmallPrint   null.String `db:"small_print"`
	YearlyPlanID string      `db:"yearly_plan_id"`
}

// ProductSchema is the db scan target.
// Description fields needs to be split into arrays by \r\n.
type ProductSchema struct {
	BaseProduct
	Description string `db:"description"`
}

// GetPlanIDs extracts the id of all plans from all products
// retrieved from DB.
func GetPlanIDs(products []ProductSchema) []string {
	var ids = make([]string, 0)

	for _, product := range products {
		ids = append(ids, product.YearlyPlanID)
	}

	return ids
}

// Product describes a product present to user
// on UI.
type Product struct {
	BaseProduct
	Description []string
	Plan        Plan
}

func ZipProductWithPlan(rows []ProductSchema, planStore GroupedPlans) []Product {
	products := make([]Product, 0)

	for _, row := range rows {

		product := Product{
			BaseProduct: row.BaseProduct,
			Description: strings.Split(row.Description, "\r\n"),
			Plan:        planStore[row.YearlyPlanID],
		}

		products = append(products, product)
	}

	return products
}
