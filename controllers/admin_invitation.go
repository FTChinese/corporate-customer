package controllers

import (
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

func NewInvitationRouter(env repository.Env, office postoffice.PostOffice) InvitationRouter {
	return InvitationRouter{
		repo: env,
		post: office,
	}
}

// Send creates an invitation for a licence and send it to a user.
// Input: {email: string, description: string, licenceId: string}
func (router InvitationRouter) Send(c echo.Context) error {
	claims := getAccountClaims(c)

	var input admin.InvitationInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// TODO: validation

	input.TeamID = claims.TeamID.String

	assignee, err := router.repo.CreateInvitation(input)
	if err != nil {
		switch err {
		case repository.ErrLicenceUnavailable:
			return &render.ValidationError{
				Message: "The licence is already taken",
				Field:   "licenceId",
				Code:    "already_taken",
			}

		case repository.ErrAlreadyMember:
			return &render.ValidationError{
				Message: "The email to accept the invitation is already a valid member",
				Field:   "membership",
				Code:    render.CodeAlreadyExists,
			}

		default:
			return render.NewDBError(err)
		}
	}

	go func() {
		// We do not use ExpandedLicence here since the
		// Assignee is still unknown until the invitation accepted.
		licence, err := router.repo.RetrieveLicence(input.LicenceID, input.TeamID)
		if err != nil {
			return
		}

		accountTeam, err := router.repo.AccountTeam(claims.Id)
		if err != nil {
			return
		}

		parcel, err := admin.ComposeInvitationLetter(assignee, licence, accountTeam)
		if err != nil {
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			logger.WithField("trace", "DeliverInvitationLetter").Error(err)
		}
	}()

	return c.NoContent(http.StatusNoContent)
}

// List shows all invitations
func (router InvitationRouter) List(c echo.Context) error {

	return nil
}

// Revoke cancels an invitation before it is accepted by user.
// If the invitation is already accepted, revoke has no effect.
// Admin should revoke a licence for this purpose.
func (router InvitationRouter) Revoke(c echo.Context) error {
	return nil
}
