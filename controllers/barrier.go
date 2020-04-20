package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// BarrierRouter defines what an admin can do before login.
type BarrierRouter struct {
	repo repository.Env
	post postoffice.PostOffice
}

// NewBarrierRouter creates a new instance of BarrierRouter.
func NewBarrierRouter(repo repository.Env, p postoffice.PostOffice) BarrierRouter {
	return BarrierRouter{
		repo: repo,
		post: p,
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
	var input admin.AccountInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	if ve := input.ValidateLogin(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	jwtAccount, err := router.repo.Login(input)
	if err != nil {
		return render.NewDBError(err)
	}

	// `200 OK`
	return c.JSON(http.StatusOK, jwtAccount)
}

// SignUp creates a new account.
// Input: {email: string, password: string}
func (router BarrierRouter) SignUp(c echo.Context) error {
	var input admin.AccountInput
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

	jwtAccount, err := router.repo.LoadPassport(signUp.ID)

	go func() {
		parcel, err := admin.ComposeVerificationLetter(jwtAccount.Account, signUp)
		if err != nil {
			return
		}

		_ = router.post.Deliver(parcel)
	}()

	// `200 OK`
	return c.JSON(http.StatusOK, jwtAccount)
}

// PostForgotPassword handles sending email to help reset password.
// Input: {email: string}
func (router BarrierRouter) PasswordResetEmail(c echo.Context) error {
	var input admin.AccountInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	// Ensure email exists in request.
	if ve := input.ValidateEmail(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Find account by email
	account, err := router.repo.AccountByEmail(input.Email)
	if err != nil {
		return render.NewDBError(err)
	}

	// Generate token.
	pr, err := input.PasswordResetter()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	// Save token
	err = router.repo.SavePasswordResetter(pr)
	if err != nil {
		return render.NewDBError(err)
	}

	// Create email content.
	parcel, err := admin.ComposePwResetLetter(account, pr)
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

	account, err := router.repo.AccountToResetPassword(token)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// ResetPassword resets password.
// Input: {token: string, password: string}
func (router BarrierRouter) ResetPassword(c echo.Context) error {
	var input admin.AccountInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidatePwReset(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountToResetPassword(input.Token)
	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.UpdatePassword(input.Credentials(account.ID))

	if err != nil {
		return render.NewDBError(err)
	}

	_ = router.repo.RemovePasswordResetToken(input.Token)

	return c.NoContent(http.StatusNoContent)
}
