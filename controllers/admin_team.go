package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TeamRouter struct {
	repo repository.Env
}

// Create creates a team for an admin.
// This is a required and first step after admin signed up successfully.
// A team is always required to perform any perform any purchase.
// After the team is created, client should refresh JWT so that
// the jwt in subsequence request will contain the team id.
// Input: {name: string, invoiceTitle?: string}
func (router TeamRouter) Create(c echo.Context) error {
	claims := getAccountClaims(c)

	var t admin.Team
	if err := c.Bind(&t); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := t.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	t = t.WithID(claims.AdminID)

	if err := router.repo.CreateTeam(t); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, t)
}

func (router TeamRouter) Load(c echo.Context) error {
	claims := getAccountClaims(c)

	t, err := router.repo.TeamByAdminID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, t)
}

// Update updates a team's name and invoice title.
// Input: {name: string, invoiceTitle?: string}
func (router TeamRouter) Update(c echo.Context) error {
	claims := getAccountClaims(c)

	var newVal admin.Team
	if err := c.Bind(&newVal); err != nil {
		return render.NewBadRequest(err.Error())
	}

	currentTeam, err := router.repo.TeamByAdminID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	if currentTeam.IsEqual(newVal) {
		return c.JSON(http.StatusOK, currentTeam)
	}

	currentTeam = currentTeam.Update(newVal)

	return c.JSON(http.StatusOK, currentTeam)
}
