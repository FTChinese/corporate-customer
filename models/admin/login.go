package admin

import (
	"github.com/FTChinese/b2b/models/validator"
	"strings"
)

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Errors   map[string]string
}

func (l *Login) Sanitize() *Login {
	l.Email = strings.TrimSpace(l.Email)
	l.Password = strings.TrimSpace(l.Password)

	return l
}

func (l *Login) Validate() bool {
	l.Errors = make(map[string]string)

	msg := validator.New("邮箱").Required().Email().Validate(l.Email)
	if msg != "" {
		l.Errors["email"] = msg
	}

	msg = validator.New("密码").Required().Validate(l.Password)
	if msg != "" {
		l.Errors["password"] = msg
	}

	return len(l.Errors) == 0
}
