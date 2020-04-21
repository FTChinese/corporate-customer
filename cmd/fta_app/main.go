package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/b2b/internal/app/b2b"
	"github.com/FTChinese/b2b/internal/app/b2b/controller"
	"github.com/FTChinese/b2b/internal/app/b2b/repository/login"
	"github.com/FTChinese/b2b/internal/app/b2b/repository/products"
	"github.com/FTChinese/b2b/internal/app/b2b/repository/setting"
	"github.com/FTChinese/b2b/internal/app/b2b/repository/subs"
	"github.com/FTChinese/b2b/internal/pkg/config"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

var (
	isProduction bool
	version      string
	build        string
	conf         config.Config
	logger       = logrus.WithField("project", "b2b").WithField("package", "main")
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
	db, err := config.NewDB(config.MustGetDBConn(conf))
	if err != nil {
		log.Fatal(err)
	}

	emailConn := config.MustGetEmailConn()
	post := postoffice.New(
		emailConn.Host,
		emailConn.Port,
		emailConn.User,
		emailConn.Pass)

	appKey := config.MustGetAppKey("web_app.b2b")

	dk := controller.NewDoorkeeper(appKey.GetJWTKey())
	subsRepo := subs.NewEnv(db)
	loginRepo := login.NewEnv(db)
	settingRepo := setting.NewEnv(db)
	productsRepo := products.NewEnv(db)

	barrierRouter := controller.NewBarrierRouter(loginRepo, post, dk)
	accountRouter := controller.NewAccountRouter(settingRepo, post, dk)
	teamRouter := controller.NewTeamRouter(subsRepo)
	productRouter := controller.NewProductRouter(productsRepo)
	orderRouter := controller.NewOrderRouter(subsRepo, productsRepo, post)
	licenceRouter := controller.NewLicenceRouter(subsRepo)
	invRouter := controller.NewInvitationRouter(subsRepo, post)
	readerRouter := controller.NewReaderRouter(subsRepo, post, dk)

	e := echo.New()
	e.Pre(middleware.AddTrailingSlash())

	e.Renderer = MustNewRenderer(conf)
	e.HTTPErrorHandler = errorHandler

	if !isProduction {
		e.Static("/css", "client/node_modules/bootstrap/dist/css")
		e.Static("/js", "client/node_modules/bootstrap.native/dist")
		e.Static("/static", "build/dev")
	}

	e.Use(b2b.DumpRequest)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())

	e.GET("/b2b/", func(context echo.Context) error {
		return context.Render(http.StatusOK, "home.html", nil)
	})

	api := e.Group("/api")
	api.POST("/login/", barrierRouter.Login)
	api.POST("/signup/", barrierRouter.SignUp)
	api.POST("/verify/:token", barrierRouter.VerifyAccount)

	pwResetGroup := api.Group("/password-reset")
	{
		// Handle resetting password
		pwResetGroup.POST("/", barrierRouter.ResetPassword)

		// Sending forgot-password email
		pwResetGroup.POST("/letter/", barrierRouter.PasswordResetEmail)

		// Verify forgot-password token.
		// If valid, redirect to /forgot-password.
		// If invalid, redirect to /forgot-password/letter to ask
		// user to enter email again.
		pwResetGroup.GET("/token/:token/", barrierRouter.VerifyPasswordToken)
	}

	accountGroup := api.Group("/account", dk.RequireLoggedIn)
	{
		accountGroup.GET("/", accountRouter.Account)
		accountGroup.GET("/jwt/", accountRouter.RefreshJWT)
		accountGroup.GET("/profile/", accountRouter.Profile)
		accountGroup.POST("/request-verification", accountRouter.RequestVerification)
		accountGroup.PATCH("/display-name", accountRouter.ChangeName)
		accountGroup.PATCH("/password", accountRouter.ChangePassword)
	}

	teamGroup := api.Group("/team", dk.RequireLoggedIn)
	{
		teamGroup.GET("/", teamRouter.Load)
		teamGroup.POST("/", teamRouter.Create)
		teamGroup.PATCH("/", teamRouter.Update)
		teamGroup.GET("/members", teamRouter.ListMembers)
		teamGroup.DELETE("/member/:id", teamRouter.DeleteMember)
	}

	productGroup := api.Group("/products", dk.RequireLoggedIn)
	{
		productGroup.GET("/", productRouter.ListProducts)
	}

	orderGroup := api.Group("/orders", dk.RequireLoggedIn)
	{
		// List orders
		orderGroup.GET("/", orderRouter.ListOrders)
		// Create orders, or renew/upgrade in bulk.
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
		invitationGroup.GET("/", invRouter.List)
		// Create invitation.
		// Also update the linked licence's status.
		invitationGroup.POST("/", invRouter.Send)
		// Revoke invitation before licence is accepted.
		// Also revert the status of a licence from invitation sent
		// back to available.
		invitationGroup.DELETE("/:id", invRouter.Revoke)
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
		readerGroup.GET("/verify-token/:token", readerRouter.VerifyInvitation)
		// Pass back data acquired from previous step
		// and get back licence data.
		readerGroup.GET("/verify-licence", readerRouter.VerifyLicence, dk.CheckInviteeClaims)
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
		readerGroup.POST("/grant", readerRouter.Grant, dk.CheckInviteeClaims)
	}

	e.Logger.Fatal(e.Start(":3100"))
}
