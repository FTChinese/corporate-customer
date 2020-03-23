package controllers

var SiteMap = struct {
	Home           string
	Login          string
	SignUp         string
	ForgotPassword string
	LogOut         string
}{
	Home:           "/",
	Login:          "/login",
	SignUp:         "/signup",
	LogOut:         "/logout",
	ForgotPassword: "/forgot-password",
}
