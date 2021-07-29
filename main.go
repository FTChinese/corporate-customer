package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/ftacademy/internal/app/b2b/controller"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/api"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

var (
	isProduction bool
	version      string
	build        string
	conf         config.Config
)

func init() {
	flag.BoolVar(&isProduction, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	config.MustSetupViper()

	conf = config.Config{
		Debug:   !isProduction,
		Version: version,
		BuiltAt: build,
	}
}

func main() {
	logger := config.MustGetLogger(isProduction)

	myDBs := db.MustNewMyDBs(isProduction)

	pm := postman.New(config.MustGetHanqiConn())

	appKey := config.MustGetAppKey("web_app.b2b")

	apiClient := api.NewSubsAPIClient(isProduction)
	dk := controller.NewDoorkeeper(appKey.GetJWTKey())

	adminRouter := controller.NewAdminRouter(myDBs, pm, dk, logger)
	subsRouter := controller.NewSubsRouter(myDBs, pm, logger)

	productRouter := controller.NewProductRouter(apiClient, logger)
	readerRouter := controller.NewReaderRouter(apiClient)

	e := echo.New()
	e.Renderer = MustNewRenderer(conf)

	if !isProduction {
		e.Static("/static", "build/public/static")
	}

	e.Pre(middleware.AddTrailingSlash())

	e.HTTPErrorHandler = errorHandler

	e.Use(controller.DumpRequest)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())

	e.GET("/corporate/*", func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "no-cache")
		c.Response().Header().Add("Cache-Control", "no-store")
		c.Response().Header().Add("Cache-Control", "must-revalidate")
		c.Response().Header().Add("Pragma", "no-cache")

		return c.Render(http.StatusOK, "b2b/home.html", nil)
	})

	api := e.Group("/api")

	b2bAPIGroup := api.Group("/b2b")

	b2bAuthGroup := b2bAPIGroup.Group("/auth")
	{
		b2bAuthGroup.POST("/login/", adminRouter.Login)
		b2bAuthGroup.POST("/signup/", adminRouter.SignUp)
		b2bAuthGroup.GET("/verify/:token/", adminRouter.VerifyEmail)

		pwResetGroup := b2bAuthGroup.Group("/password-reset")
		{
			// Handle resetting password
			pwResetGroup.POST("/", adminRouter.ResetPassword)

			// Sending forgot-password email
			pwResetGroup.POST("/letter/", adminRouter.ForgotPassword)

			// Verify forgot-password token.
			// If valid, redirect to /forgot-password.
			// If invalid, redirect to /forgot-password/letter to ask
			// user to enter email again.
			pwResetGroup.GET("/token/:token/", adminRouter.VerifyResetToken)
		}
	}

	b2bAccountGroup := b2bAPIGroup.Group("/account", dk.RequireLoggedIn)
	{
		//b2bAccountGroup.GET("/", accountRouter.Account)
		b2bAccountGroup.GET("/jwt/", adminRouter.RefreshJWT)
		b2bAccountGroup.POST("/request-verification/", adminRouter.RequestVerification)
		b2bAccountGroup.PATCH("/display-name/", adminRouter.ChangeName)
		b2bAccountGroup.PATCH("/password/", adminRouter.ChangePassword)
	}

	b2bTeamGroup := b2bAPIGroup.Group("/team/", dk.RequireLoggedIn)
	{
		b2bTeamGroup.GET("/", adminRouter.LoadTeam)
		b2bTeamGroup.POST("/", adminRouter.CreateTeam)
		b2bTeamGroup.PATCH("/", adminRouter.UpdateTeam)
	}

	productGroup := b2bAPIGroup.Group("/paywall", dk.RequireLoggedIn)
	{
		productGroup.GET("/", productRouter.Paywall)
	}

	orderGroup := b2bAPIGroup.Group("/orders", dk.RequireLoggedIn)
	{
		// List orders
		orderGroup.GET("/", subsRouter.ListOrders)
		// CreateTeam orders, or renew/upgrade in bulk.
		orderGroup.POST("/", subsRouter.CreateOrders)
		orderGroup.GET("/:id/", subsRouter.LoadOrder)
	}

	b2bLicenceGroup := b2bAPIGroup.Group("/licences", dk.RequireLoggedIn)
	{
		// List licences
		b2bLicenceGroup.GET("/", subsRouter.ListLicence)
		b2bLicenceGroup.GET("/:id/", subsRouter.LoadLicence)
		// Revoked a licence
		b2bLicenceGroup.POST("/:id/revoke/", subsRouter.RevokeLicence)
	}

	b2bInvitationGroup := b2bAPIGroup.Group("/invitations", dk.RequireLoggedIn)
	{
		// List invitations
		b2bInvitationGroup.GET("/", subsRouter.ListInvitations)
		// CreateTeam invitation.
		// Also update the linked licence's status.
		b2bInvitationGroup.POST("/", subsRouter.CreateInvitation)
		// Revoked invitation before licence is accepted.
		// Also revert the status of a licence from invitation sent
		// back to available.
		b2bInvitationGroup.POST("/:id/revoke/", subsRouter.RevokeInvitation)
	}

	// Steps to accept an invitation:
	// 1. Open token url and the token is valid;
	// 2. Use email to find user account (If account not found, go to signup);
	// 3. Get account data and find out if membership already exists
	// 4. Grant licence
	b2bGrantGroup := b2bAPIGroup.Group("/grant-licence")
	{
		// Verify the invitation is valid. Cache the invitation for a short period
		// so that the next step won't hit db.
		b2bGrantGroup.GET("/verify-invitation/:token/", subsRouter.VerifyInvitation)
		b2bGrantGroup.POST("/signup/", readerRouter.SignUp)
		// Grant licence to user:
		// 1. Retrieve invitation again;
		// 2. Use invitation email to get reader account and verify it again.
		// 3. Lock invitation row, lock licence row, lock membership row if exists.
		// 4. Set invitation being used; link licence to reader id; backup existing
		// membership if exists; upsert membership.
		// 5. Sent email to reader and admin about the result.
		b2bGrantGroup.POST("/grant/", subsRouter.GrantLicence)
	}

	e.Logger.Fatal(e.Start(":4000"))
}
