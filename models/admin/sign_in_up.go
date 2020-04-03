package admin

import (
	"github.com/FTChinese/b2b/models/form"
	"github.com/FTChinese/go-rest/rand"
	"github.com/google/uuid"
)

// Login is used to verify user's credentials.
type Login struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

// NewLogin creates a new instance of  Login from input data.
func NewLogin(a form.AccountForm) Login {
	return Login{
		Email:    a.Email,
		Password: a.Password,
	}
}

// Verifier holds email verification data.
type Verifier struct {
	ID    string `db:"admin_id"`
	Token string `db:"token"`
}

// NewVerifier creates a new verification token for a new user
// if id is empty, or regenerate a token for existing user
// if id is provided.
func NewVerifier(id string) (Verifier, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return Verifier{}, err
	}

	if id == "" {
		id = uuid.New().String()
	}
	return Verifier{
		ID:    id,
		Token: token,
	}, nil
}

// SignUp is user's sign up data.
type SignUp struct {
	Verifier
	Login
}

// SignUp creates a new instance of SignUp from input data.
func NewSignUp(a form.AccountForm) (SignUp, error) {

	v, err := NewVerifier("")
	if err != nil {
		return SignUp{}, err
	}

	return SignUp{
		Verifier: v,
		Login: Login{
			Email:    a.Email,
			Password: a.Password,
		},
	}, nil
}

// PasswordResetter holds the token to identify password resetting owner.
type PasswordResetter struct {
	Email string `db:"email"`
	Token string `db:"token"`
}

func NewPasswordResetter(email string) (PasswordResetter, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return PasswordResetter{}, err
	}
	return PasswordResetter{
		Token: token,
		Email: email,
	}, nil
}

// Credentials is used to update user's password.
type Credentials struct {
	ID       string `db:"admin_id"`
	Password string `db:"password"`
}
