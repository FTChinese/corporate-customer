package views

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/widget"
)

func NewResetLetterForm(a admin.AccountForm) widget.Form {
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
