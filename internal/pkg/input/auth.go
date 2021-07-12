package input

import (
	"github.com/FTChinese/ftacademy/pkg/validator"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"strings"
)

type Credentials struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func (c *Credentials) Validate() *render.ValidationError {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)

	ve := validator.
		New("email").
		Required().
		MaxLen(64).
		Email().
		Validate(c.Email)

	if ve != nil {
		return ve
	}

	return validator.
		New("password").
		Required().
		MaxLen(64).
		Validate(c.Password)
}

type SignupParams struct {
	Credentials
	SourceURL string `json:"sourceUrl"` // URL to verify email.
}

// ForgotPasswordParams is used to create a password reset session.
type ForgotPasswordParams struct {
	Email     string      `json:"email"`
	SourceURL null.String `json:"sourceUrl"`
}

func (f ForgotPasswordParams) Validate() *render.ValidationError {
	f.Email = strings.TrimSpace(f.Email)

	return validator.EnsureEmail(f.Email)
}

// PasswordResetParams contains the data used to reset
type PasswordResetParams struct {
	Token    string `json:"token"`    // identify this session
	Password string `json:"password"` // the new password user submitted
}

func (i *PasswordResetParams) Validate() *render.ValidationError {
	i.Token = strings.TrimSpace(i.Token)
	i.Password = strings.TrimSpace(i.Password)

	ve := validator.
		New("token").
		Required().
		Validate(i.Token)

	if ve != nil {
		return ve
	}

	return validator.EnsurePassword(i.Password)
}
