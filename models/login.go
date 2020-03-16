package models

import "strings"

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Errors   map[string]string
}

func (l *Login) Validate() bool {
	l.Email = strings.TrimSpace(l.Email)
	l.Password = strings.TrimSpace(l.Password)

	l.Errors = make(map[string]string)

	if l.Email == "" {
		l.Errors["Email"] = "User name cannot be empty"
	}

	if l.Password == "" {
		l.Errors["Password"] = "Password cannot be empty"
	}

	if len(l.Password) < 6 {
		l.Errors["Password"] = "Password too short"
	}

	return len(l.Errors) == 0
}
