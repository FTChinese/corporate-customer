package controller

import (
	"github.com/FTChinese/b2b/internal/app/b2b/model"
	"github.com/FTChinese/b2b/internal/app/b2b/repository/subs"
	"github.com/FTChinese/b2b/internal/app/b2b/stmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TeamRouter struct {
	repo subs.Env
}

func NewTeamRouter(repo subs.Env) TeamRouter {
	return TeamRouter{
		repo: repo,
	}
}

// Create creates a team for an admin.
// This is a required and first step after admin signed up successfully.
// A team is always required to perform any perform any purchase.
// After the team is created, client should refresh JWT so that
// the jwt in subsequence request will contain the team id.
// Input: {name: string, invoiceTitle?: string}
func (router TeamRouter) Create(c echo.Context) error {
	claims := getPassportClaims(c)

	var t model.Team
	if err := c.Bind(&t); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := t.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	t = t.BuildOn(claims.AdminID)

	if err := router.repo.CreateTeam(t); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, t)
}

func (router TeamRouter) Load(c echo.Context) error {
	claims := getPassportClaims(c)

	t, err := router.repo.TeamByAdminID(claims.AdminID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, t)
}

// Update updates a team's name and invoice title.
// Input: {name: string, invoiceTitle?: string}
func (router TeamRouter) Update(c echo.Context) error {
	claims := getPassportClaims(c)

	var newVal model.Team
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

func (router TeamRouter) ListMembers(c echo.Context) error {
	claims := getPassportClaims(c)
	var page gorest.Pagination
	if err := c.Bind(page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	countCh, listCh := router.repo.AsyncCountStaff(stmt.TeamByID), router.repo.AsyncListStaff(claims.TeamID.String, page)

	countResult, listResult := <-countCh, <-listCh
	if listResult.Err != nil {
		return render.NewDBError(listResult.Err)
	}

	listResult.Total = countResult.Total
	return c.JSON(http.StatusOK, listResult)
}

// Delete removes an assignee.
func (router TeamRouter) DeleteMember(c echo.Context) error {
	claims := getPassportClaims(c)

	memberID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	err = router.repo.DeleteStaffer(model.Staffer{
		ID:     memberID,
		TeamID: claims.TeamID.String,
	})

	if err != nil {
		return render.NewDBError(err)
	}

	return nil
}
