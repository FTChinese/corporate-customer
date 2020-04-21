package view

import (
	"github.com/FTChinese/b2b/internal/app/b2b/model"
	"github.com/FTChinese/b2b/pkg/widget"
)

func NewResetLetterForm(a model.AccountInput) widget.Form {
	return widget.Form{
		Disabled: false,
		Method:   widget.MethodPost,
		Action:   "",
		Fields: []widget.FormControl{
			{
				Label:       "邮箱",
				ID:          "email",
				Type:        widget.ControlTypeEmail,
				Name:        "email",
				Value:       a.Email,
				Placeholder: "admin@example.org",
				Required:    true,
			},
		},
		SubmitBtn: widget.PrimaryBlockBtn.
			SetName("发送邮件").
			SetDisabledText("正在发送..."),
		CancelBtn: widget.Link{},
		DeleteBtn: widget.Link{},
	}
}
