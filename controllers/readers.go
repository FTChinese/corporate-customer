package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/reader"
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReadersRouter struct {
	repo repository.Env
	post postoffice.PostOffice
}

func NewReaderRouter(env repository.Env, post postoffice.PostOffice) ReadersRouter {
	return ReadersRouter{
		repo: env,
		post: post,
	}
}

// VerifyInvitation verifies an invitation by token.
// Status code:
// 404 if the invitation does not exist, expired, or already used.
func (router ReadersRouter) VerifyInvitation(c echo.Context) error {
	token := c.Param("token")

	inv, err := router.repo.FindInvitationByToken(token)
	// Not found error.
	if err != nil {
		return render.NewDBError(err)
	}

	// If invitation is expires or already accepted or revoked.
	if !inv.IsValid() {
		return render.NewNotFound("the invitation is expired, revoked or already used")
	}

	bearer, err := admin.NewInvitationBearer(admin.NewInviteeClaims(inv))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.JSON(http.StatusOK, bearer)
}

// VerifyLicence checks whether a licence being invited
// is available to grant.
// Status code:
// 404 if the licence is not found for this invitation. or cannot be granted.
func (router ReadersRouter) VerifyLicence(c echo.Context) error {
	claims := getInviteeClaims(c)

	lic, err := router.repo.FindInvitedLicence(claims)
	if err != nil {
		return render.NewDBError(err)
	}

	if !lic.CanBeGranted() {
		return render.NewNotFound("the licence is already granted, or not tailed for this invitation")
	}

	return c.NoContent(http.StatusNoContent)
}

// FindAccount gets an invited reader's account.
// Status code:
// 404 if reader is not signed up yet.
// 422 if this reader already has a valid membership.
func (router ReadersRouter) FindAccount(c echo.Context) error {
	claims := getInviteeClaims(c)

	account, err := router.repo.FindReader(claims.Email)
	if err != nil {
		// For no found, ask user to signup.
		return render.NewDBError(err)
	}

	if account.HasMembership() {
		return render.NewUnprocessable(&render.ValidationError{
			Message: "Cannot have multiple subscription since you already have a valid one",
			Field:   "membership",
			Code:    render.CodeAlreadyExists,
		})
	}

	claims.FtcID = account.FtcID.String
	bearer, err := admin.NewInvitationBearer(claims)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.JSON(http.StatusOK, bearer)
}

// SignUp create a new account for an invited reader.
// Input: {email: string, password: string}
// Status code:
// 400 bad request.
// 422 if input data is not valid.
// Returns a new JWT token containing ftc id if everything
// went well.
func (router ReadersRouter) SignUp(c echo.Context) error {
	claims := getInviteeClaims(c)

	var input admin.AccountInput
	if err := c.Bind(input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateLogin(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	signUp, err := reader.NewSignUp(input)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	if err := router.repo.CreateReader(signUp); err != nil {
		return render.NewDBError(err)
	}

	// Add the missing ftc id for a team member.
	go func() {
		_ = router.repo.UpdateTeamMember(signUp.TeamMember(claims.TeamID))
	}()

	// Add the missing ftc id for new reader.
	claims.FtcID = signUp.ID
	bearer, err := admin.NewInvitationBearer(claims)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}
	return c.JSON(http.StatusOK, bearer)
}

// Grant links a licence to a reader invited to accept it.
// Status code:
// 400 if invitation is not found, or is invalid,
// or licence is not found or cannot be granted,
// or account is not found.
// 403 Forbidden if reader already has valid membership.
func (router ReadersRouter) Grant(c echo.Context) error {
	claims := getInviteeClaims(c)

	if claims.FtcID == "" {
		return render.NewNotFound("user not found")
	}

	il, err := router.repo.GrantLicence(claims)
	if err != nil {
		switch err {
		case repository.ErrInvalidInvitation:
			return render.NewNotFound(err.Error())

		case repository.ErrLicenceTaken:
			return render.NewNotFound(err.Error())

		default:
			return render.NewDBError(err)
		}
	}

	// Send a notification letter to admin.
	go func() {
		pp, err := router.repo.PassportByAdminID(claims.TeamID)
		if err != nil {
			return
		}

		parcel, err := admin.ComposeLicenceGranted(il, pp)
		if err != nil {
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			logger.WithField("trace", "ReaderRouter").Error(err)
		}
	}()

	return c.NoContent(http.StatusNoContent)
}
