package admin

import (
	"github.com/FTChinese/b2b/models/validator"
	"github.com/FTChinese/go-rest/rand"
	"github.com/FTChinese/go-rest/render"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"strings"
)

func GenerateToken() (string, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return "", err
	}

	return token, nil
}

// AccountInput is used to parse account-related
// request body.
// Different forms might use a combination of
// different fields, resulting to duplicate
// fields if we create a data type for each
// form.
//			    | Request		   |		DB
// -----------------------------------------------------------------
// Login        | Email + Password | Email + Password
// Sign Up      | Email + Password | Email + Password + ID + Token
// Reset letter | Email            | N/A
// Reset pw     | Password + Token | ID + Password
// Update pw    | Password         | ID + Password
// Display name | DisplayName      | ID + DisplayName
type AccountInput struct {
	ID          string      `db:"admin_id"`
	Email       string      `json:"email" db:"email"`
	Password    string      `json:"password" db:"password"`
	Token       string      `json:"token" db:"token"`
	OldPassword string      `json:"oldPassword" db:"old_password"`
	DisplayName null.String `json:"displayName" db:"display_name"`
	IsSignUp    bool        // Used a verification letter template to determine the greeting message
}

// NewVerifier regenerates a verification token for a new user.
func NewVerifier(id string) (AccountInput, error) {
	token, err := GenerateToken()
	if err != nil {
		return AccountInput{}, err
	}

	return AccountInput{
		ID:    id,
		Token: token,
	}, nil
}

func (a *AccountInput) ValidateEmail() *render.ValidationError {
	a.Email = strings.TrimSpace(a.Email)

	return validator.New("email").Required().Email().Validate(a.Email)
}

func (a *AccountInput) ValidateDisplayName() *render.ValidationError {
	name := strings.TrimSpace(a.DisplayName.String)
	a.DisplayName = null.NewString(name, name != "")

	return validator.New("displayName").Max(64).Validate(name)
}

// ValidatePassword validates password fields.
// Used when user is updating password.
func (a *AccountInput) ValidatePassword() *render.ValidationError {
	a.Password = strings.TrimSpace(a.Password)

	ve := validator.New("password").Required().Validate(a.Password)
	if ve != nil {
		return ve
	}

	return nil
}

func (a *AccountInput) ValidatePasswordUpdate() *render.ValidationError {
	a.OldPassword = strings.TrimSpace(a.OldPassword)

	ve := validator.New("oldPassword").Required().Validate(a.OldPassword)
	if ve != nil {
		return ve
	}

	return a.ValidatePassword()
}

// ValidateLogin performs validation for login
// or signup.
func (a *AccountInput) ValidateLogin() *render.ValidationError {

	if ve := a.ValidateEmail(); ve != nil {
		return ve
	}

	if ve := a.ValidatePassword(); ve != nil {
		return ve
	}

	return nil
}

// SignUp creates data fields required to created a new account.
// It need ID and Token in addition to user input fields Email and Password.
func (a AccountInput) SignUp() (AccountInput, error) {
	token, err := GenerateToken()
	if err != nil {
		return a, err
	}

	a.ID = uuid.New().String()
	a.Token = token
	a.IsSignUp = true
	return a, nil
}

// Letter produce a Letter type.
// This is used to unify the parameter passed to build
// verification letter since a verification letter might
// be needed under two situations: upon SignUp or
// user' explicit request.
// Call this on the returned value of SignUp or NewVerifier.
func (a AccountInput) VerificationLetter() Letter {
	return Letter{
		URL:      "https://www.ftacademy.cn/b2b/verify/" + a.Token,
		IsSignUp: a.IsSignUp,
	}
}

// ValidatePwReset validates resetting password and its associated token.
func (a *AccountInput) ValidatePwReset() *render.ValidationError {
	a.Token = strings.TrimSpace(a.Token)

	ve := validator.New("token").Required().Validate(a.Token)
	if ve != nil {
		return ve
	}

	return a.ValidatePassword()
}

// PasswordResetter produces the data to save a password resetting row.
func (a AccountInput) PasswordResetter() (AccountInput, error) {
	token, err := GenerateToken()
	if err != nil {
		return a, err
	}

	return AccountInput{
		Email: a.Email,
		Token: token,
	}, nil
}

// ResettingLetter produces the data to compile template.
// Call this from the returned value of PasswordResetter.
func (a AccountInput) ResettingLetter() Letter {
	return Letter{
		URL: "https://www.ftacademy.cn/b2b/password-reset/token/" + a.Token,
	}
}

// Credentials produces a new AccountInput with ID and Password
// fields set to updated a user's password.
func (a AccountInput) Credentials(id string) AccountInput {
	return AccountInput{
		ID:       id,
		Password: a.Password,
	}
}
