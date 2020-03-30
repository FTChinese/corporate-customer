package admin

import (
	"github.com/FTChinese/b2b/models/validator"
	"github.com/FTChinese/go-rest/rand"
	"github.com/google/uuid"
	"strings"
)

// AccountForm is used to parse form data related to
// administrator's account.
// Different forms use different fields combination.
// Login: Email + Password
// Sign Up: Email + Password + ConfirmPassword
// Request password reset letter: Email
// Reset password: Password + ConfirmPassword
// Update password: OldPassword + Password + ConfirmPassword
type AccountForm struct {
	Email           string `form:"email"`
	DisplayName     string `form:"displayName"`
	OldPassword     string `form:"oldPassword"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirmPassword"`
	Errors          map[string]string
}

func (a *AccountForm) validateEmail() {
	a.Email = strings.TrimSpace(a.Email)

	msg := validator.New("邮箱").Required().Email().Validate(a.Email)
	if msg != "" {
		a.Errors["email"] = msg
	}
}

func (a *AccountForm) validatePassword() {
	a.Password = strings.TrimSpace(a.Password)

	msg := validator.New("密码").Required().Validate(a.Password)
	if msg != "" {
		a.Errors["password"] = msg
	}
}

func (a *AccountForm) validateConfirmation() {
	a.ConfirmPassword = strings.TrimSpace(a.ConfirmPassword)

	msg := validator.New("确认密码").Required().Validate(a.ConfirmPassword)
	if msg != "" {
		a.Errors["confirmPassword"] = msg
		return
	}

	if a.Password != a.ConfirmPassword {
		a.Errors["confirmPassword"] = "两次输入的密码不一致"
	}
}

// ValidateLogin performs validation for login.
func (a *AccountForm) ValidateLogin() bool {

	a.Errors = make(map[string]string)

	a.validateEmail()
	a.validatePassword()

	return len(a.Errors) == 0
}

// ValidateSignUp performs validation for sign-up.
func (a *AccountForm) ValidateSignUp() bool {
	a.Errors = make(map[string]string)

	a.validateEmail()
	a.validatePassword()
	a.validateConfirmation()

	return len(a.Errors) == 0
}

// ValidateEmail performs validation for password
// resetting email.
func (a *AccountForm) ValidateEmail() bool {
	a.Errors = make(map[string]string)

	a.validateEmail()

	return len(a.Errors) == 0
}

func (a *AccountForm) ValidatePasswordReset() bool {
	a.Errors = make(map[string]string)

	a.validatePassword()
	a.validateConfirmation()

	return len(a.Errors) == 0
}

// Login creates Login type from input data
func (a AccountForm) NewLogin() Login {
	return Login{
		Email:    a.Email,
		Password: a.Password,
	}
}

// SignUp create SignUp type from input data.
func (a AccountForm) NewSignUp() (SignUp, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return SignUp{}, err
	}
	return SignUp{
		ID:           uuid.New().String(),
		Email:        a.Email,
		Password:     a.Password,
		Verification: token,
	}, nil
}
