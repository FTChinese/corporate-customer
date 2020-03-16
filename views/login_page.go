package views

import (
	"github.com/FTChinese/corporate-customer/models"
	"github.com/flosch/pongo2"
)

func BuildLoginForm(value models.Login) Form {
	return Form{
		Disabled: false,
		Action:   "",
		Fields: []FormField{
			{
				Label:       "邮箱",
				ID:          "email",
				Type:        InputTypeEmail,
				Name:        "email",
				Value:       value.Email,
				Placeholder: "admin@example.org",
				Required:    true,
			},
			{
				Label:    "密码",
				ID:       "password",
				Type:     InputTypePassword,
				Name:     "password",
				Value:    "",
				Required: true,
			},
		},
		SubmitBtn: Button{
			DisableWith: "正在登录",
			Text:        "登录",
		},
		CancelBtn: Link{},
		DeleteBtn: Link{},
	}
}

func BuildLoginPage(value models.Login) pongo2.Context {
	return pongo2.Context{
		"form": BuildLoginForm(value),
	}
}
