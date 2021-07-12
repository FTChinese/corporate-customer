package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateTeam creates a team for an admin.
// This is a required and first step after admin signed up successfully.
// A team is always required to perform any perform any purchase.
// After the team is created, client should refresh JWT so that
// the jwt in subsequence request will contain the team id.
// Input: {name: string, invoiceTitle?: string}
func (router AdminRouter) CreateTeam(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getPassportClaims(c)

	var params input.TeamParams
	if err := c.Bind(&params); err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		sugar.Error(ve)
		return render.NewUnprocessable(ve)
	}

	team := admin.NewTeam(claims.AdminID, params)

	if err := router.repo.CreateTeam(team); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, team)
}

func (router AdminRouter) LoadTeam(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getPassportClaims(c)

	t, err := router.repo.LoadTeam(claims.AdminID)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, t)
}

// UpdateTeam updates a team's name and invoice title.
// Input: {orgName: string, invoiceTitle?: string}
func (router AdminRouter) UpdateTeam(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	claims := getPassportClaims(c)

	var newVal input.TeamParams
	if err := c.Bind(&newVal); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := newVal.Validate(); ve != nil {
		sugar.Error(ve)
		return render.NewUnprocessable(ve)
	}

	currentTeam, err := router.repo.LoadTeam(claims.AdminID)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	if currentTeam.IsEqual(newVal) {
		return c.JSON(http.StatusOK, currentTeam)
	}

	teamUpdated := currentTeam.Update(newVal)

	err = router.repo.UpdateTeam(teamUpdated)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, teamUpdated)
}
