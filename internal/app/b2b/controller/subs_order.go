package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateOrders creates orders an org purchased.
// input: input.ShoppingCart.
func (router SubsRouter) CreateOrders(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getPassportClaims(c)

	var cart checkout.ShoppingCart
	if err := c.Bind(&cart); err != nil {
		return render.NewBadRequest(err.Error())
	}

	schema := checkout.NewOrderSchemaBuilder(cart, claims).
		Build()
	err := router.repo.CreateOrder(schema)
	if err != nil {
		return render.NewDBError(err)
	}

	go func() {
		profile, err := router.repo.LoadB2BAdminProfile(claims.AdminID)
		if err != nil {
			sugar.Error(err)
			return
		}

		parcel, err := letter.OrderCreatedParcel(profile, schema.OrderRow)
		if err != nil {
			sugar.Error(err)
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			sugar.Error(err)
		}
	}()

	return c.JSON(http.StatusOK, schema.OrderRow)
}

func (router SubsRouter) ListOrders(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	list, err := router.repo.ListOrders(
		claims.TeamID.String,
		page)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, list)
}

func (router SubsRouter) LoadOrder(c echo.Context) error {
	claims := getPassportClaims(c)

	id := c.Param("id")

	o, err := router.repo.RetrieveOrder(admin.AccessRight{
		RowID:  id,
		TeamID: claims.TeamID.String,
	})

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, o)
}
