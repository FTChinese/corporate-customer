package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/letter"
	subs2 "github.com/FTChinese/ftacademy/internal/repository/subs"
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
		case subs2.ErrInvalidInvitation:
			return render.NewNotFound(err.Error())

		case subs2.ErrLicenceTaken:
			return render.NewNotFound(err.Error())

		default:
			return render.NewDBError(err)
		}
	}

	// Send a notification letter to admin.
	go func() {
		// Find the admin by licence creator id.
		profile, err := router.repo.AdminProfile(result.Licence.CreatorID)
		if err != nil {
			sugar.Infof("Error retreiving admin profile %s of licence %s", result.Licence.CreatorID, result.Licence.ID)
			sugar.Error(err)
			return
		}

		parcel, err := letter.LicenceGrantedParcel(result.Licence, profile)
		if err != nil {
			sugar.Infof("Error creating granted parcel for licence %s", result.Licence.ID)
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

	// TODO: send email to user.

	return c.JSON(http.StatusOK, result)
}
