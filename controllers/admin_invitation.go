package controllers

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InvitationRouter struct {
	repo repository.Env
	post postoffice.PostOffice
}

// Send creates an invitation for a licence and send it to a user.
// Input: {email: string, description: string, licenceId: string}
func (router InvitationRouter) Send(c echo.Context) error {
	claims := getAccountClaims(c)

	var input admin.InvitationInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// Find the licence by input.LicenceID
	licence, err := router.repo.LoadLicence(input.LicenceID, claims.TeamID.String)
	if err != nil {
		return render.NewDBError(err)
	}

	if !licence.IsAvailable() {
		return &render.ValidationError{
			Message: "The licence is already taken",
			Field:   "licenceId",
			Code:    "already_taken",
		}
	}

	// Find the user by email
	r, err := router.repo.FindReader(input.Email)
	// Ignore not found error.
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return render.NewDBError(err)
	}
	// If this use already has a valid membership.
	// A non-existing user is always treated as expired.
	if !r.Membership.IsExpired() {
		return &render.ValidationError{
			Message: "The email to accept the invitation is already a valid member",
			Field:   "membership",
			Code:    render.CodeAlreadyExists,
		}
	}

	// Now licence is available, user does not have a valid
	// membership, you can grant the licence to this user.

	// Compose letter.

	// Send letter

	// You also save this assignee under current admin.

	return c.NoContent(http.StatusNoContent)
}

// List shows all invitations
func (router InvitationRouter) List(c echo.Context) error {

	return nil
}

// Revoke cancels an invitation before it is accepted by user.
// If the invitation is already accepted, revoke has no effect.
// Admin should revoke a licence for this purpose.
func (router InvitationRouter) Revoke() error {
	return nil
}
