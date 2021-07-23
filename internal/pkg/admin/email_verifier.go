package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"time"
)

// EmailVerifier holds the data used to create an email
// to verify the authenticity of reader's email address.
// For sign-up, the request body contains `email` and `sourceUrl`, and IsSignUp should be true
// For manually requesting verification, request body contains only `sourceUrl`. The email comes from database.
// Client should sent the SourceURL when asking for a verification letter so that
// the client could run under any host.
type EmailVerifier struct {
	Token        string      `db:"token"`
	Email        string      `db:"email"`
	ExpireInDays int64       `db:"expire_in_days"`
	CreatedUTC   chrono.Time `db:"created_utc"`
}

// NewEmailVerifier creates a verification token for an email.
// It sets a default url to build the verification link.
func NewEmailVerifier(email string) (EmailVerifier, error) {
	token, err := gorest.RandomHex(32)

	if err != nil {
		return EmailVerifier{}, err
	}

	return EmailVerifier{
		Token:        token,
		Email:        email,
		ExpireInDays: 3,
		CreatedUTC:   chrono.TimeNow(),
	}, nil
}

func (v EmailVerifier) IsExpired() bool {
	return v.CreatedUTC.AddDate(0, 0, int(v.ExpireInDays)).Before(time.Now())
}

// BuildURL creates a verification link.
func (v EmailVerifier) BuildURL() string {
	return pkg.B2BVerifyAdminURL(v.Token)
}
