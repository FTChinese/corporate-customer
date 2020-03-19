package admin

import (
	"github.com/FTChinese/b2b/models/validator"
	"strings"
)

// SignUp is user's sign up data.
type SignUp struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirmPassword"`
	Errors          map[string]string
}

// Sanitize trims leading and trailing spaces.
func (s *SignUp) Sanitize() *SignUp {
	s.Email = strings.TrimSpace(s.Email)
	s.Password = strings.TrimSpace(s.Password)
	s.ConfirmPassword = strings.TrimSpace(s.ConfirmPassword)

	return s
}

// Validate checks if all fields are valid.
func (s *SignUp) Validate() bool {
	s.Errors = make(map[string]string)

	msg := validator.New("邮箱").Required().Email().Validate(s.Email)
	if msg != "" {
		s.Errors["email"] = msg
	}

	msg = validator.New("密码").Required().Validate(s.Password)
	if msg != "" {
		s.Errors["password"] = msg
	}

	msg = validator.New("确认密码").Required().Validate(s.ConfirmPassword)
	if msg != "" {
		s.Errors["confirmPassword"] = msg
	}

	return len(s.Errors) == 0
}

// ConfirmationMatched checks whether the confirmation password
// matches the password.
func (s *SignUp) ConfirmationMatched() bool {
	return s.Password == s.ConfirmPassword
}
