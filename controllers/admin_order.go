package controllers

import (
	"github.com/FTChinese/b2b/repository"
	"github.com/labstack/echo/v4"
)

type OrderRouter struct {
	repo repository.Env
}

// CreateOrders creates orders an org purchased.
// Client should specify which plans are being subscribed,
// and how many copies.
// Input:
// [
//     {planId: "string", quantity: number},
//     {planId: "string", quantity: number}
// ]
// At most there should be two plans: a standard and a premium.
func (router OrderRouter) CreateOrders(c echo.Context) error {
	//claims := getAccountClaims(c)

	return nil
}
