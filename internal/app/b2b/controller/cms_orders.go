package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router CMSRouter) ListOrders(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	teamID := c.QueryParam("team_id")
	status := c.QueryParam("status")
	s, _ := checkout.ParseStatus(status)
	filter := checkout.OrderFilter{
		TeamID: teamID,
		Status: s,
	}

	list, err := router.repo.ListOrders(filter, page)

	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, list)
}

func (router CMSRouter) LoadOrder(c echo.Context) error {
	orderID := c.Param("id")

	o, err := router.repo.LoadDetailedOrder(orderID)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, o)
}

// ConfirmPayment confirms payment of an order.
// Input:
// amountPaid: number;
// description?: string;
// paymentMethod: string;
// transactionId: string;
// priceOffPerCopy: [{ orderItemId: string, priceOff: number }]
func (router CMSRouter) ConfirmPayment(c echo.Context) error {
	orderID := c.Param("id")

	// Get operator name
	// Change order payment fields
	// Update price_off_per_copy of order_item
	// TODO: what's next?

	return c.JSON(http.StatusOK, orderID)
}
