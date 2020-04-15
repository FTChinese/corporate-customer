package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OrderRouter struct {
	repo repository.Env
}

// CreateOrders creates orders an org purchased.
// Client should specify which plans are being subscribed,
// and how many copies.
// Input:
// [
//     {planId: "string", quantity: number, cycleCount: number },
//     {planId: "string", quantity: number, cycleCount: number }
// ]
// At most there should be two plans: a standard and a premium.
func (router OrderRouter) CreateOrders(c echo.Context) error {
	claims := getAccountClaims(c)

	var cartItems []admin.CartItem
	if err := c.Bind(&cartItems); err != nil {
		return render.NewBadRequest(err.Error())
	}

	plans, err := router.repo.PlansInSet(admin.GetCartPlanIDs(cartItems))

	if err != nil {
		return render.NewDBError(err)
	}

	cart, ve := admin.NewCart(cartItems, plans)
	if ve != nil {
		return render.NewUnprocessable(ve)
	}

	orders := cart.BuildOrders(claims.TeamID.String)

	err = router.repo.CreateOrders(orders)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router OrderRouter) ListOrders(c echo.Context) error {
	return nil
}
