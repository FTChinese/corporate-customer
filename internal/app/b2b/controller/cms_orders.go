package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
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

	o, err := router.repo.LoadOrder(orderID)

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
// offers: [{
// 		copies: number;
//		kind: 'create' | 'renew';
//		price: Price;
//		priceOffPerCopy: number;
// }]
func (router CMSRouter) ConfirmPayment(c echo.Context) error {
	orderID := c.Param("id")

	var params input.OrderPaidParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// Retrieve order.
	order, err := router.repo.LoadOrder(orderID)
	if err != nil {
		return render.NewDBError(err)
	}

	if order.IsFinal() {
		return render.NewBadRequest("Order already paid")
	}

	// Update order status
	order = order.ChangeStatus(checkout.StatusProcessing)
	err = router.repo.UpdateOrderStatus(order)
	if err != nil {
		return render.NewDBError(err)
	}

	// Save payment
	payResult := checkout.NewOrderPid(order.ID, params)
	err = router.repo.SavePaymentResult(payResult)
	if err != nil {
		return render.NewDBError(err)
	}

	// TODO: run without blocking.
	err = router.repo.ConfirmPayment(order)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, payResult)
}
