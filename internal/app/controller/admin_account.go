package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// RefreshJWT updates jwt token.
func (router AdminRouter) RefreshJWT(c echo.Context) error {
	claims := getAdminClaims(c)

	baseAccount, err := router.repo.BaseAccountByID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	bearer, err := router.guard.CreatePassport(baseAccount)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, bearer)
}

// RequestVerification sends user a new verification letter.
// Used must be logged in.
// Input: none.
// Status codes:
// 404 - The account for this user is not found.
// 500 - Token generation failed or DB error.
func (router AdminRouter) RequestVerification(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getAdminClaims(c)
	var params input.ReqEmailVrfParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// Find the account
	account, err := router.repo.BaseAccountByID(claims.AdminID)
	// 404 Not Found
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	re := router.sendEmailVerification(account)

	if re != nil {
		return re
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeName updates display name.
// Input: {displayName: string}.
// StatusCodes:
// 400 - If request body cannot be parsed.
// Client should refresh JWT after success.
func (router AdminRouter) ChangeName(c echo.Context) error {
	claims := getAdminClaims(c)

	var params input.NameUpdateParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.BaseAccountByID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	if account.DisplayName.String == params.DisplayName {
		return c.NoContent(http.StatusNoContent)
	}

	updated := account.UpdateName(params.DisplayName)

	if err := router.repo.UpdateName(updated); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, updated)
}

// ChangePassword updates password.
// Input: {oldPassword: string, password: string}
// Status codes:
// 400 - Request cannot be parsed.
// 422 -
// {
// 	message: "",
//  error: { field: "password", code: "invalid" }
// }
// 404 - Account not found.
// 401 - Unauthorized, meaning old password is not correct.
// 500 - DB error.
// 204 - Success.
func (router AdminRouter) ChangePassword(c echo.Context) error {
	claims := getAdminClaims(c)

	var params input.PasswordUpdateParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	params.ID = claims.AdminID

	// Verify old password.
	authResult, err := router.repo.VerifyPassword(params)
	if err != nil {
		return render.NewDBError(err)
	}

	if !authResult.PasswordMatched {
		return render.NewUnauthorized("Wrong old password")
	}

	// Change to new password.
	if err := router.repo.UpdatePassword(params); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
