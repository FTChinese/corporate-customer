package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/api"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

type ReadersRouter struct {
	apiClient api.Client
}

func NewReaderRouter(client api.Client) ReadersRouter {
	return ReadersRouter{
		apiClient: client,
	}
}

// SignUp create a new account for an invited reader.
// Input: {email: string, password: string}
// Status code:
// 400 bad request.
// 422 if input data is not valid.
// Returns a new JWT token containing ftc id if everything
// went well.
func (router ReadersRouter) SignUp(c echo.Context) error {

	var params input.SignupParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	params.SourceURL = config.B2BReaderVerification

	resp, err := router.apiClient.ReaderSignup(params)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReadersRouter) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")

	resp, err := router.apiClient.VerifySignup(token)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
