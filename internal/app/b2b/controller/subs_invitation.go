package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateInvitation creates an invitation for a licence and send it to a user.
// Input:
// email: string,
// description: string,
// licenceId: string
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
	params.TeamID = claims.TeamID.String

	lic, err := router.repo.CreateInvitation(params, claims.AdminID)
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
		assignee, err := router.repo.FindAssignee(params.Email)
		if err != nil {
			sugar.Error(err)
			return
		}
		adminProfile, err := router.repo.AdminProfile(claims.AdminID)
		if err != nil {
			sugar.Error(err)
			return
		}

		parcel, err := letter.InvitationParcel(
			assignee,
			lic,
			adminProfile)
		if err != nil {
			sugar.Error(err)
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			sugar.Error(err)
		}
	}()

	return c.JSON(http.StatusOK, licence.Licence{
		BaseLicence: lic,
		Assignee:    licence.AssigneeJSON{},
	})
}

// RevokeInvitation revokes the invitation of a licence
// before the licence is granted.
// If the invitation is already accepted, revoke has no effect.
// Admin should revoke a licence for this purpose.
func (router SubsRouter) RevokeInvitation(c echo.Context) error {
	invID := c.Param("id") // the invitation id
	claims := getPassportClaims(c)

	result, err := router.repo.RevokeInvitation(invID, claims.TeamID.String)
	// TODO: handle different errors
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, result)
}

// ListInvitations shows all invitations
// Query: ?page=1&per_page=10
func (router SubsRouter) ListInvitations(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	invs, err := router.repo.ListInvitations(claims.TeamID.String, page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, invs)
}

func (router SubsRouter) VerifyInvitation(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	token := c.QueryParam("token")

	inv, err := router.repo.InvitationByToken(token)
	if err != nil {
		sugar.Error(err)
		// 404 Not Found
		return render.NewDBError(err)
	}

	// Ensure the invitation is usable.
	if !inv.IsAcceptable() {
		return render.NewBadRequest("invitation is either used or expired")
	}

	lic, err := router.repo.LoadLicence(admin.AccessRight{
		RowID:  inv.LicenceID,
		TeamID: inv.TeamID,
	})
	if err != nil {
		sugar.Error(err)
		// 404 Not Found
		return render.NewDBError(err)
	}
	// TODO: remove this or keep it? We could leave it to the client to determine whether the licence is available.
	if !lic.IsAvailable() {
		sugar.Error(err)
		return render.NewBadRequest("Licence is not available to be granted")
	}

	// Find user by invitation email.
	// If user is not found, it only indicates user is not
	// signed-up.
	assignee, err := router.repo.FindAssignee(inv.Email)
	// 404 won't happen here
	if err != nil {
		return render.NewDBError(err)
	}
	if assignee.IsZero() {
		c.JSON(http.StatusOK, licence.InvitationVerified{
			Licence:    lic,
			Assignee:   assignee,
			Membership: reader.Membership{},
		})
	}

	m, err := router.repo.RetrieveMembership(assignee.FtcID.String)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, licence.InvitationVerified{
		Licence:    lic,
		Assignee:   assignee,
		Membership: m,
	})
}
