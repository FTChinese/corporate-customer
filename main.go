package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/b2b/controllers"
	"github.com/FTChinese/b2b/database"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	config       Config
	logger       = logrus.WithField("project", "superyard").WithField("package", "main")
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

	config = Config{
		Debug:   !isProduction,
		Version: version,
		BuiltAt: build,
		Year:    0,
	}
}

func main() {
	db, err := database.New(config.MustGetDBConn("mysql.master"))
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	emailConn := MustGetEmailConn()
	post := postoffice.NewPostman(
		emailConn.Host,
		emailConn.Port,
		emailConn.User,
		emailConn.Pass)

	signInRouter := controllers.NewSignInRouter(db, post)
	readersRouter := controllers.NewReadersRouter(db, post)

	e := echo.New()
	e.Renderer = MustNewRenderer(config)
	e.HTTPErrorHandler = errorHandler

	if !isProduction {
		e.Static("/css", "client/node_modules/bootstrap/dist/css")
		e.Static("/js", "client/node_modules/bootstrap.native/dist")
		e.Static("/static", "build/dev")
	}

	e.Use(middleware.Logger())
	e.Use(session.Middleware(
		sessions.NewCookieStore(
			[]byte(MustGetSessionKey()),
		),
	))
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())

	e.GET("/", func(context echo.Context) error {
		return context.Render(http.StatusOK, "home.html", nil)
	}, controllers.NotLoggedIn)

	e.GET("/login", signInRouter.GetLogin, controllers.AlreadyLoggedIn)
	e.POST("/login", signInRouter.PostLogin)

	e.GET("/logout", signInRouter.LogOut)

	pwResetGroup := e.Group("/password-reset")
	pwResetGroup.GET("/", signInRouter.GetResetPassword)
	pwResetGroup.POST("/", signInRouter.PostResetPassword)

	pwResetGroup.GET("/letter", signInRouter.GetForgotPassword)
	pwResetGroup.POST("/letter", signInRouter.PostResetPassword)

	pwResetGroup.GET("/token/:token", signInRouter.VerifyPasswordToken)

	e.GET("/readers", readersRouter.GetUserList, controllers.NotLoggedIn)
	//readersGroup := e.Group("/readers")

	e.Logger.Fatal(e.Start(":3100"))
}
