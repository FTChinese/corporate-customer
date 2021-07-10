package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/login"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// BarrierRouter defines what an admin can do before login.
type BarrierRouter struct {
	keeper Doorkeeper
	repo   login.Env
	post   postman.Postman
}

// NewBarrierRouter creates a new instance of BarrierRouter.
func NewBarrierRouter(repo login.Env, p postman.Postman, dk Doorkeeper) BarrierRouter {
	return BarrierRouter{
		keeper: dk,
		repo:   repo,
		post:   p,
	}
}

// Login verifies email + password.
// Input: {email: string, password: string}
// Status code:
// 400 Bad Request if request body cannot be parsed.
// 404 Not Found if email or password, or both, are incorrect.
// 422 Unprocessable
// * if email is missing:
// { error: { field: "email", code: "missing_field"}}
// * if email is invalid:
// {error: {field: "email", code: "invalid"}}
// * if password is missing:
// {error: { field: "password", code: "missing_field"}}
func (router BarrierRouter) Login(c echo.Context) error {
	var input model.AccountInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	if ve := input.ValidateLogin(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	passport, err := router.repo.Login(input)
	if err != nil {
		return render.NewDBError(err)
	}

	jwtBearer, err := passport.Bearer(router.keeper.signingKey)
	if err != nil {
		return render.NewInternalError(err.Error())
	}
	// `200 OK`
	return c.JSON(http.StatusOK, jwtBearer)
}

// SignUp creates a new account.
// Input: {email: string, password: string}
func (router BarrierRouter) SignUp(c echo.Context) error {
	var input model.AccountInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	if ve := input.ValidateLogin(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	signUp, err := input.SignUp()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	err = router.repo.SignUp(signUp)
	if err != nil {
		return render.NewDBError(err)
	}

	// Convert to Passport type.
	pp := signUp.Passport()

	go func() {
		parcel, err := model.ComposeVerificationLetter(pp.Account, signUp)
		if err != nil {
			return
		}

		_ = router.post.Deliver(parcel)
	}()

	jwtBearer, err := pp.Bearer(router.keeper.signingKey)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	// `200 OK`
	return c.JSON(http.StatusOK, jwtBearer)
}

// PostForgotPassword handles sending email to help reset password.
// Input: {email: string}
func (router BarrierRouter) PasswordResetEmail(c echo.Context) error {
	var input model.AccountInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	// Ensure email exists in request.
	if ve := input.ValidateEmail(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Find account by email
	account, err := router.repo.FindPwResetAccount(input.Email)
	if err != nil {
		return render.NewDBError(err)
	}

	// Generate token.
	pr, err := input.PasswordResetter()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	// Save token
	err = router.repo.SavePwResetToken(pr)
	if err != nil {
		return render.NewDBError(err)
	}

	// Create email content.
	parcel, err := model.ComposePwResetLetter(account, pr)
	if err != nil {
		return err
	}

	// Send email in background.
	go func() {
		_ = router.post.Deliver(parcel)
	}()

	return c.NoContent(http.StatusNoContent)
}

// VerifyPasswordToken when user clicked link in an email.
// No input. Get the token from path parameter.
func (router BarrierRouter) VerifyPasswordToken(c echo.Context) error {
	token := c.Param("token")

	account, err := router.repo.FindAccountByPwToken(token)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// ResetPassword resets password.
// Input: {token: string, password: string}
func (router BarrierRouter) ResetPassword(c echo.Context) error {
	var input model.AccountInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidatePwReset(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.FindAccountByPwToken(input.Token)
	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.ResetPassword(input.Credentials(account.ID))

	if err != nil {
		return render.NewDBError(err)
	}

	_ = router.repo.RemovePwResetToken(input.Token)

	return c.NoContent(http.StatusNoContent)
}

// VerifyEmail finds the account associated with a token
// and set is as verified if found.
// Status codes:
// 404 - Account not found for this token.
// 204 - Verified successfully.
func (router BarrierRouter) VerifyAccount(c echo.Context) error {
	token := c.Param("token")

	account, err := router.repo.VerifyingAccount(token)
	if err != nil {
		return render.NewDBError(err)
	}
	// If it is already verified, return immediately.
	if account.Verified {
		return c.NoContent(http.StatusNoContent)
	}

	account.Verified = true
	err = router.repo.SetAccountVerified(account)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
