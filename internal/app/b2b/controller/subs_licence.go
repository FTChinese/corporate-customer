package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subsrepo"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router SubsRouter) LoadLicence(c echo.Context) error {
	licID := c.Param("licID")
	// TODO: ensure teamId actually exist before hitting this endpoint.
	claims := getPassportClaims(c)

	lic, err := router.repo.LoadLicence(admin.AccessRight{
		RowID:  licID,
		TeamID: claims.TeamID.String,
	})
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, lic)
}

func (router SubsRouter) ListLicence(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	licences, err := router.repo.ListLicence(claims.TeamID.String, page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, licences)
}

// GrantLicence links a licence to a reader invited to accept it.
// Input:
// * licenceId: string;
// * ftcId : string;
//
// Status code:
// 400 if invitation is not found, or is invalid,
// or licence is not found or cannot be granted,
// or account is not found.
// 403 Forbidden if reader already has valid membership.
func (router SubsRouter) GrantLicence(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	var params input.GrantParams
	if err := c.Bind(&params); err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	assignee, err := router.repo.RetrieveAssignee(params.FtcID)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	result, err := router.repo.GrantLicence(admin.AccessRight{
		RowID:  params.LicenceID,
		TeamID: params.TeamID,
	}, assignee)

	// TODO: handle various errors.
	if err != nil {
		sugar.Error(err)
		switch err {
		case subsrepo.ErrInvalidInvitation:
			return render.NewNotFound(err.Error())

		case subsrepo.ErrLicenceTaken:
			return render.NewNotFound(err.Error())

		default:
			return render.NewDBError(err)
		}
	}

	// Send a notification letter to admin.
	go func() {
		err := router.repo.SaveVersionedLicence(result.LicenceVersion)
		if err != nil {
			sugar.Error()
		}

		err = router.repo.ArchiveMembership(result.MemberModified.MembershipVersion)
		if err != nil {
			sugar.Error(err)
		}

		// TODO: send email to this user.

		// Find the admin by licence creator id.
		profile, err := router.repo.LoadB2BAdminProfile(result.LicenceVersion.PostChange.AdminID)
		if err != nil {
			sugar.Infof("Error retreiving admin profile %s of licence %s", result.LicenceVersion.PostChange.AdminID, result.LicenceVersion.PostChange.ID)
			sugar.Error(err)
			return
		}

		parcel, err := letter.LicenceGrantedParcel(result.LicenceVersion.PostChange.Licence, assignee, profile)
		if err != nil {
			sugar.Infof("Error creating granted parcel for licence %s", result.LicenceVersion.PostChange.ID)
			sugar.Error(err)
			return
		}

		err = router.post.Deliver(parcel)
		if err != nil {
			sugar.Error(err)
		}
	}()

	// Send back user's membership.
	return c.JSON(http.StatusOK, result)
}

// RevokeLicence unlinks a reader from a licence.
func (router SubsRouter) RevokeLicence(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	id := c.Param("id")

	claims := getPassportClaims(c)

	result, err := router.repo.RevokeLicence(admin.AccessRight{
		RowID:  id,
		TeamID: claims.TeamID.String,
	})

	// TODO: handle various error response
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	go func() {
		err := router.repo.SaveVersionedLicence(result.LicenceVersion)
		if err != nil {
			sugar.Error()
		}

		err = router.repo.ArchiveMembership(result.MembershipVersioned)
		if err != nil {
			sugar.Error(err)
		}

		// TODO: if membership has addon, send a request to API to re-enable it.
		// TODO: send email to this user.
	}()

	return c.JSON(http.StatusOK, result)
}
