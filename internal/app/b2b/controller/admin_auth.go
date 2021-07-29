package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SignUp creates a new account.
// Input: {email: string, password: string}
func (router AdminRouter) SignUp(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	var params input.SignupParams
	if err := c.Bind(&params); err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		sugar.Error(ve)
		return render.NewUnprocessable(ve)
	}

	adminAccount := admin.NewAccount(params)

	err := router.repo.SignUp(adminAccount)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	go func() {
		_ = router.sendEmailVerification(adminAccount.BaseAccount)
	}()

	jwtBearer, err := admin.NewPassport(adminAccount.BaseAccount, router.keeper.signingKey)
	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	// `200 OK`
	return c.JSON(http.StatusOK, jwtBearer)
}

func (router AdminRouter) sendEmailVerification(baseAccount admin.BaseAccount) *render.ResponseError {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	verifier, err := admin.NewEmailVerifier(baseAccount.Email)
	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	// Save new token.
	err = router.repo.SaveEmailVerifier(verifier)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	// Generate letter content
	parcel, err := letter.VerificationParcel(baseAccount, verifier)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	err = router.post.Deliver(parcel)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return nil
}

// VerifyEmail finds the account associated with a token
// and set is as verified if found.
// Status codes:
// 404 - BaseAccount not found for this token.
// 204 - Verified successfully.
func (router AdminRouter) VerifyEmail(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	token := c.Param("token")

	vrf, err := router.repo.RetrieveEmailVerifier(token)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}
	if vrf.IsExpired() {
		return render.NewNotFound("Token expired")
	}

	baseAccount, err := router.repo.BaseAccountByEmail(vrf.Email)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}
	if baseAccount.Verified {
		return c.NoContent(http.StatusNoContent)
	}

	err = router.repo.EmailVerified(baseAccount.ID)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
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
func (router AdminRouter) Login(c echo.Context) error {
	var params input.Credentials
	if err := c.Bind(&params); err != nil {
		return err
	}

	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	authResult, err := router.repo.Authenticate(params)
	if err != nil {
		return render.NewDBError(err)
	}
	if !authResult.PasswordMatched {
		return render.NewForbidden("Incorrect credentials")
	}

	baseAccount, err := router.repo.BaseAccountByID(authResult.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	jwtBearer, err := admin.NewPassport(baseAccount, router.keeper.signingKey)
	if err != nil {
		return render.NewInternalError(err.Error())
	}
	// `200 OK`
	return c.JSON(http.StatusOK, jwtBearer)
}

// ForgotPassword sends email to help user reset password.
// Input: {email: string, sourceUrl?: string}
func (router AdminRouter) ForgotPassword(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	var params input.ForgotPasswordParams
	if err := c.Bind(&params); err != nil {
		sugar.Error(err)
		return err
	}

	// Ensure email exists in request.
	if ve := params.Validate(); ve != nil {
		sugar.Error(ve)
		return render.NewUnprocessable(ve)
	}

	// Find account by email
	account, err := router.repo.BaseAccountByEmail(params.Email)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	// Generate token.
	session, err := admin.NewPwResetSession(params)
	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	// Save token
	err = router.repo.SavePwResetSession(session)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	// CreateTeam email content.
	parcel, err := letter.PasswordResetParcel(account, session)
	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	err = router.post.Deliver(parcel)
	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// VerifyResetToken when user clicked link in an email.
// No input. Get the token from path parameter.
func (router AdminRouter) VerifyResetToken(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	token := c.Param("token")

	session, err := router.repo.PwResetSession(token)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}
	if session.IsUsed || session.IsExpired() {
		return render.NewNotFound("token already used or expired")
	}

	return c.JSON(http.StatusOK, session)
}

// ResetPassword resets password.
// Input: {token: string, password: string}
func (router AdminRouter) ResetPassword(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	var params input.PasswordResetParams
	if err := c.Bind(&params); err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		sugar.Error(ve)
		return render.NewUnprocessable(ve)
	}

	session, err := router.repo.PwResetSession(params.Token)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	baseAccount, err := router.repo.BaseAccountByEmail(session.Email)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	err = router.repo.UpdatePassword(input.PasswordUpdateParams{
		ID:  baseAccount.ID,
		New: params.Password,
	})

	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	go func() {
		_ = router.repo.DisablePasswordReset(session.WithUsed())
	}()

	return c.NoContent(http.StatusNoContent)
}
