package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AdminAccountRouter struct {
	repo repository.Env
	post postoffice.PostOffice
}

func NewAccountRouter(repo repository.Env, post postoffice.PostOffice) AdminAccountRouter {
	return AdminAccountRouter{
		repo: repo,
		post: post,
	}
}

// RefreshJWT updates jwt token.
func (router AdminAccountRouter) RefreshJWT(c echo.Context) error {
	claims := getAccountClaims(c)

	jwtAccount, err := router.repo.JWTAccount(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, jwtAccount)
}

// Account sends user's account data.
func (router AdminAccountRouter) Account(c echo.Context) error {
	claims := getAccountClaims(c)

	account, err := router.repo.AccountByID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// Profile sends user's profile.
// Status codes:
// 404 - Account not found
// 200 - with Profile as body.
func (router AdminAccountRouter) Profile(c echo.Context) error {
	claims := getAccountClaims(c)

	profile, err := router.repo.AdminProfile(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, profile)
}

// RequestVerification sends user a new verification letter.
// Used must be logged in.
// Input: none.
// Status codes:
// 404 - The account for this user is not found.
// 500 - Token generation failed or DB error.
func (router AdminAccountRouter) RequestVerification(c echo.Context) error {
	claims := getAccountClaims(c)

	// Find the account
	account, err := router.repo.AccountByID(claims.AdminID)
	// 404 Not Found
	if err != nil {
		return render.NewDBError(err)
	}

	// Generate new verification token.
	verifier, err := admin.NewVerifier(claims.AdminID)
	// 500
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	// Save new token.
	err = router.repo.RegenerateVerifier(verifier)
	// 500
	if err != nil {
		return render.NewDBError(err)
	}

	// Generate letter content
	parcel, err := account.VerificationLetter(verifier.VerificationLetter())
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	// Send letter
	go func() {
		_ = router.post.Deliver(parcel)
	}()

	return c.NoContent(http.StatusNoContent)
}

// VerifyEmail finds the account associated with a token
// and set is as verified if found.
// Status codes:
// 404 - Account not found for this token.
// 204 - Verified successfully.
func (router AdminAccountRouter) VerifyEmail(c echo.Context) error {
	token := c.Param("token")

	account, err := router.repo.AccountByVerifier(token)
	if err != nil {
		return render.NewDBError(err)
	}
	// If it is already verified, return immediately.
	if account.Verified {
		return c.NoContent(http.StatusNoContent)
	}

	account.Verified = true
	err = router.repo.SetEmailVerified(account)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeName updates display name.
// Input: {displayName: string}.
// StatusCodes:
// 400 - If request body cannot be parsed.
func (router AdminAccountRouter) ChangeName(c echo.Context) error {
	claims := getAccountClaims(c)

	var input admin.AccountInput
	if err := c.Bind(input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateDisplayName(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	if account.DisplayName == input.DisplayName {
		return c.NoContent(http.StatusNoContent)
	}

	account.DisplayName = input.DisplayName
	if err := router.repo.UpdateName(account); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
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
func (router AdminAccountRouter) ChangePassword(c echo.Context) error {
	claims := getAccountClaims(c)

	var input admin.AccountInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidatePasswordUpdate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	input.ID = claims.AdminID
	// Verify old password.
	matched, err := router.repo.PasswordMatched(input)
	if err != nil {
		return render.NewDBError(err)
	}

	if !matched {
		return render.NewUnauthorized("Wrong old password")
	}

	// Change to new password.
	if err := router.repo.UpdatePassword(input); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
