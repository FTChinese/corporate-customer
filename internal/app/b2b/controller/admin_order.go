package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/products"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	"github.com/FTChinese/ftacademy/pkg/postman"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OrderRouter struct {
	repo         subs.Env
	productsRepo products.Env
	post         postman.Postman
}

func NewOrderRouter(env subs.Env, prodRepo products.Env, p postman.Postman) OrderRouter {
	return OrderRouter{
		repo:         env,
		productsRepo: prodRepo,
		post:         p,
	}
}

// CreateOrders creates orders an org purchased.
// Client should specify which plans are being subscribed,
// and how many copies.
// TODO: design the data structure of a cart.
// Input:
// [
//  {planId: string, quantity: number, cycleCount: number, kind: "create" },
//  {planId: string, cycleCount: number, kind: "renew", licences: string[] }
//	{planId: string, cycleCount: number, kind: "upgrade", licences: string[] }
// ]
// At most there should be two plans: a standard and a premium.
func (router OrderRouter) CreateOrders(c echo.Context) error {
	claims := getPassportClaims(c)

	var cartItems []model.CartItem
	if err := c.Bind(&cartItems); err != nil {
		return render.NewBadRequest(err.Error())
	}

	plans, err := router.productsRepo.PlansInSet(model.GetCartPlanIDs(cartItems))

	if err != nil {
		return render.NewDBError(err)
	}

	cart, ve := model.NewCart(cartItems, plans)
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
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	listCh, countCh := router.repo.AsyncListOrders(claims.TeamID.String, page), router.repo.AsyncCountOrder(stmt.TeamByID)

	listResult, countResult := <-listCh, <-countCh
	if listResult.Err != nil {
		return render.NewDBError(listResult.Err)
	}

	listResult.Total = countResult.Total
	return c.JSON(http.StatusOK, listResult)
}
