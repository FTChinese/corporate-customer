package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/ftacademy/internal/app/b2b"
	"github.com/FTChinese/ftacademy/internal/app/b2b/controller"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/products"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var (
	isProduction bool
	version      string
	build        string
	conf         config.Config
	logger       = logrus.WithField("project", "ftacademy").WithField("package", "main")
)

func init() {
	flag.BoolVar(&isProduction, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		os.Exit(1)
	}

	conf = config.Config{
		Debug:   !isProduction,
		Version: version,
		BuiltAt: build,
	}
}

func main() {
	logger := config.MustGetLogger(isProduction)
	// TODO: use read/write/delete dbs
	myDB := db.MustNewMySQL(config.MustMySQLWriteConn(isProduction))
	myDBs := db.MustNewMyDBs(isProduction)

	pm := postman.New(config.MustGetHanqiConn())

	appKey := config.MustGetAppKey("web_app.b2b")

	dk := controller.NewDoorkeeper(appKey.GetJWTKey())
	// TODO: deprecate
	subsRepo := subs.NewEnv(myDBs, logger)
	productsRepo := products.NewEnv(myDB)

	adminRouter := controller.NewAdminRouter(myDBs, pm, dk, logger)
	subsRouter := controller.NewSubsRouter(myDBs, pm, logger)
	productRouter := controller.NewProductRouter(productsRepo)
	// TODO: deprecate
	orderRouter := controller.NewOrderRouter(subsRepo, productsRepo, pm)
	licenceRouter := controller.NewLicenceRouter(subsRepo)
	readerRouter := controller.NewReaderRouter(subsRepo, pm, dk)

	e := echo.New()
	e.Renderer = MustNewRenderer(conf)

	if !isProduction {
		e.Static("/static", "build/public/static")
	}

	e.Pre(middleware.AddTrailingSlash())

	e.HTTPErrorHandler = errorHandler

	e.Use(b2b.DumpRequest)
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

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login/", adminRouter.Login)
		authGroup.POST("/signup/", adminRouter.SignUp)
		authGroup.GET("/verify/:token", adminRouter.VerifyEmail)

		pwResetGroup := authGroup.Group("/password-reset")
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

	accountGroup := api.Group("/account", dk.RequireLoggedIn)
	{
		//accountGroup.GET("/", accountRouter.Account)
		accountGroup.GET("/jwt/", adminRouter.RefreshJWT)
		accountGroup.POST("/request-verification", adminRouter.RequestVerification)
		accountGroup.PATCH("/display-name", adminRouter.ChangeName)
		accountGroup.PATCH("/password", adminRouter.ChangePassword)
	}

	teamGroup := api.Group("/team", dk.RequireLoggedIn)
	{
		teamGroup.GET("/", adminRouter.LoadTeam)
		teamGroup.POST("/", adminRouter.CreateTeam)
		teamGroup.PATCH("/", adminRouter.UpdateTeam)
	}

	// TODO: use subscription api.
	productGroup := api.Group("/products", dk.RequireLoggedIn)
	{
		productGroup.GET("/", productRouter.ListProducts)
	}

	// TODO: delete
	orderGroup := api.Group("/orders", dk.RequireLoggedIn)
	{
		// List orders
		orderGroup.GET("/", orderRouter.ListOrders)
		// CreateTeam orders, or renew/upgrade in bulk.
		orderGroup.POST("/", orderRouter.CreateOrders)
	}

	licenceGroup := api.Group("/licences", dk.RequireLoggedIn)
	{
		// List licences
		licenceGroup.GET("/", licenceRouter.ListLicence)
		// Renew/upgrade a licence
		licenceGroup.PATCH("/:id", licenceRouter.UpdateLicence)
		// Revoke a licence
		licenceGroup.DELETE("/:id", licenceRouter.RevokeLicence)
	}

	invitationGroup := api.Group("/invitations", dk.RequireLoggedIn)
	{
		// List invitations
		invitationGroup.GET("/", subsRouter.ListInvitations)
		// CreateTeam invitation.
		// Also update the linked licence's status.
		invitationGroup.POST("/", subsRouter.CreateInvitation)
		// Revoke invitation before licence is accepted.
		// Also revert the status of a licence from invitation sent
		// back to available.
		invitationGroup.DELETE("/:id", subsRouter.RevokeInvitation)
	}

	// Steps to accept an invitation:
	// 1. Open token url and the token is valid;
	// 2. Use email to find user account (If account not found, go to signup);
	// 3. Get account data and find out if membership already exists
	// 4. Grant licence
	readerGroup := api.Group("/accept-invitation")
	{
		// Verify the invitation is valid. Cache the invitation for a short period
		// so that the next step won't hit db.
		readerGroup.GET("/verify/:token", readerRouter.VerifyInvitation)
		// Pass back data acquired from previous step
		// and get back licence data.
		readerGroup.GET("/licence", readerRouter.VerifyLicence, dk.CheckInviteeClaims)
		// Pass back data acquired from previous step
		// and get reader account.
		// If response is not found, go to signup.
		readerGroup.GET("/account", readerRouter.FindAccount, dk.CheckInviteeClaims)
		readerGroup.POST("/signup", readerRouter.SignUp, dk.CheckInviteeClaims)
		// Grant licence to user:
		// 1. Retrieve invitation again;
		// 2. Use invitation email to get reader account and verify it again.
		// 3. Lock invitation row, lock licence row, lock membership row if exists.
		// 4. Set invitation being used; link licence to reader id; backup existing
		// membership if exists; upsert membership.
		// 5. Sent email to reader and admin about the result.
		readerGroup.POST("/grant", subsRouter.GrantLicence, dk.CheckInviteeClaims)
	}

	e.Logger.Fatal(e.Start(":4000"))
}
