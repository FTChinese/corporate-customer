package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
)

type Login struct {
	Email    string
	Password string
}

// SignUp is user's sign up data.
type SignUp struct {
	ID           string `db:"admin_id"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Verification string `db:"vrf_token"`
}

type PasswordResetter struct {
	Token string `db:"token"`
	Email string `db:"email"`
}

func NewPasswordReseter(email string) (PasswordResetter, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return PasswordResetter{}, err
	}
	return PasswordResetter{
		Token: token,
		Email: email,
	}, nil
}

type PasswordResettingAccount struct {
	Account                // The account whose password will be reset
	ExpiresIn  int64       `db:"expires_in"`  // Token duration.
	CreatedUTC chrono.Time `db:"created_utc"` // Token creation time
}

// Credentials is used to update user's password.
type Credentials struct {
	ID       string `db:"admin_id"`
	Password string `db:"password"`
}
