package products

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	plan2 "github.com/FTChinese/ftacademy/pkg/plan"
)

func (env Env) LoadProducts() ([]plan2.Product, error) {
	// TODO: find in cache.

	productRows, err := env.retrieveProducts()
	if err != nil {
		return nil, err
	}

	planIDs := plan2.GetProductsPlanIDs(productRows)

	groupedPlans, err := env.PlansInSet(planIDs)
	if err != nil {
		return nil, err
	}

	products := plan2.ZipProductWithPlan(productRows, groupedPlans)

	// TODO: cache the final products.

	return products, nil
}

func (env Env) retrieveProducts() ([]plan2.ProductSchema, error) {
	var rows []plan2.ProductSchema
	if err := env.db.Select(&rows, stmt.Products); err != nil {
		return nil, err
	}

	return rows, nil
}
