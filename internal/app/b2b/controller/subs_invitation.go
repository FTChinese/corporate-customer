package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateInvitation creates an invitation for a licence and send it to a user.
// Input: {email: string, description: string, licenceId: string}
func (router SubsRouter) CreateInvitation(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getPassportClaims(c)

	var params input.InvitationParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	inv, err := model2.NewInvitation(params, claims)
	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	invitedLicence, err := router.repo.CreateInvitation(inv)
	if err != nil {
		switch err {
		case subs.ErrLicenceUnavailable:
			return &render.ValidationError{
				Message: "The licence is already taken",
				Field:   "licence",
				Code:    "already_taken",
			}

		case subs.ErrInviteeMismatch:
			return &render.ValidationError{
				Message: err.Error(),
				Field:   "invitee",
				Code:    render.CodeAlreadyExists,
			}

		case subs.ErrAlreadyMember:
			return &render.ValidationError{
				Message: "The email to accept the invitation is already a valid member",
				Field:   "membership",
				Code:    render.CodeAlreadyExists,
			}

		default:
			return render.NewDBError(err)
		}
	}

	// Send invitation letter
	go func() {

		adminProfile, err := router.repo.AdminProfile(claims.AdminID)
		if err != nil {
			sugar.Error(err)
			return
		}

		parcel, err := letter.InvitationParcel(invitedLicence, adminProfile)
		if err != nil {
			sugar.Error()
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			sugar.Error(err)
		}
	}()

	// Add the invitee to team member
	go func() {
		_ = router.repo.SaveStaffer(invitedLicence.Assignee.TeamMember(claims.TeamID.String))
	}()

	return c.NoContent(http.StatusNoContent)
}

// RevokeInvitation cancels an invitation before it is accepted by user.
// If the invitation is already accepted, revoke has no effect.
// Admin should revoke a licence for this purpose.
func (router SubsRouter) RevokeInvitation(c echo.Context) error {
	invID := c.Param("id") // the invitation id
	claims := getPassportClaims(c)

	err := router.repo.RevokeInvitation(invID, claims.TeamID.String)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListInvitations shows all invitations
// Query: ?page=1&per_page=10
func (router SubsRouter) ListInvitations(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	countCh, listCh := router.repo.AsyncCountInvitation(claims.TeamID.String), router.repo.AsyncListInvitations(claims.TeamID.String, page)

	countResult, listResult := <-countCh, <-listCh
	if listResult.Err != nil {
		return render.NewDBError(listResult.Err)
	}

	listResult.Total = countResult.Total
	return c.JSON(http.StatusOK, listResult)
}
