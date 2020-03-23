package views

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/widget"
)

func NewLoginForm(value admin.Login) widget.Form {
	return widget.Form{
		Disabled: false,
		Action:   "",
		Fields: []widget.FormControl{
			{
				Label:       "邮箱",
				ID:          "email",
				Type:        widget.ControlTypeEmail,
				Name:        "email",
				Value:       value.Email,
				Placeholder: "admin@example.org",
				Required:    true,
			},
			{
				Label:    "密码",
				ID:       "password",
				Type:     widget.ControlTypePassword,
				Name:     "password",
				Value:    "",
				Required: true,
			},
		},
		SubmitBtn: widget.Button{
			DisableWith: "正在登录...",
			Text:        "登录",
		},
		CancelBtn: widget.Link{},
		DeleteBtn: widget.Link{},
	}
}
