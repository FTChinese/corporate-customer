package controller

import (
	"database/sql"
	"github.com/FTChinese/ftacademy/internal/app/repository/subsrepo"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/validator"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// FindMembership searches if an email has membership.
func (router SubsRouter) FindMembership(c echo.Context) error {
	email := strings.TrimSpace(c.QueryParam("email"))
	ve := validator.New("email").Email().Validate(email)
	if ve != nil {
		return render.NewUnprocessable(ve)
	}

	a, err := router.repo.FindAssignee(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, reader.Membership{})
		}

		return render.NewDBError(err)
	}

	m, err := router.repo.RetrieveMembership(a.FtcID.String)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, m)
}

// CreateInvitation creates an invitation for a licence and send it to a user.
// Input:
// email: string,
// description: string,
// licenceId: string
// Returns a licence.ExpandedLicence instance with its LatestInvitation
// field populated with the invitation created here.
func (router SubsRouter) CreateInvitation(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getAdminClaims(c)

	var params input.InvitationParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Create invitation and get the update licence.
	lic, err := router.repo.CreateInvitation(params, claims)
	if err != nil {
		switch err {
		case subsrepo.ErrLicenceUnavailable:
			return &render.ValidationError{
				Message: "The licence is already taken",
				Field:   "licence",
				Code:    "already_taken",
			}

		case subsrepo.ErrInviteeMismatch:
			return &render.ValidationError{
				Message: err.Error(),
				Field:   "invitee",
				Code:    render.CodeAlreadyExists,
			}

		case subsrepo.ErrAlreadyMember:
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
		sugar.Info("Start sending invitation email")
		assignee, err := router.repo.FindAssignee(params.Email)
		if err != nil {
			sugar.Error(err)
			return
		}
		// Find admin so that we could tell user who send the invitation.
		sugar.Info("Start retrieving b2b admin profile")
		adminProfile, err := router.repo.LoadB2BAdminProfile(claims.AdminID)
		if err != nil {
			sugar.Error(err)
			return
		}

		parcel, err := letter.InvitationParcel(
			assignee,
			lic,
			adminProfile)

		sugar.Infof("Invitation parcel %v", parcel)

		if err != nil {
			sugar.Error(err)
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			sugar.Error(err)
		} else {
			sugar.Infof("Invitation letter sent to %s", assignee.Email.String)
		}
	}()

	return c.JSON(http.StatusOK, licence.ExpandedLicence{
		Licence:  lic,
		Assignee: licence.AssigneeJSON{}, // Assignee field should be empty after invitation is created.
	})
}

// RevokeInvitation revokes the invitation of a licence
// before the licence is granted.
// If the invitation is already accepted, revoke has no effect.
// Admin should revoke a licence for this purpose.
func (router SubsRouter) RevokeInvitation(c echo.Context) error {
	invID := c.Param("id") // the invitation id
	claims := getAdminClaims(c)

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
	claims := getAdminClaims(c)

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

	token := c.Param("token")

	inv, err := router.repo.InvitationByToken(token)
	if err != nil {
		sugar.Error(err)
		// 404 Not Found
		return render.NewDBError(err)
	}

	// Ensure the invitation is usable.
	if !inv.IsAcceptable() {
		sugar.Infof("Invitation %s is cannot be granted", inv.ID)
		return render.NewBadRequest("invitation is either used or expired")
	}

	lic, err := router.repo.LoadLicence(admin.AccessRight{
		RowID:  inv.LicenceID,
		TeamID: inv.TeamID,
	})
	if err != nil {
		sugar.Infof("Error retrieve licence %s", inv.LicenceID)
		sugar.Error(err)
		// 404 Not Found
		return render.NewDBError(err)
	}
	// TODO: remove this or keep it? We could leave it to the client to determine whether the licence is available.
	if !lic.IsAvailable() {
		sugar.Infof("ExpandedLicence %s is not available to be granted", lic.ID)
		sugar.Error(err)
		return render.NewBadRequest("ExpandedLicence is not available to be granted")
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
		return c.JSON(http.StatusOK, licence.InvitationVerified{
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
