package admin

import (
	"fmt"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"time"
)

// PwResetSession hold the token to allow resetting password.
// When user is requesting password reset,
// email is sent to server to create the data; an optional
// sourceUrl should be provided by client is user is on
// desktop so that we could create a clickable link in letter.
// If user is using mobile apps, only email is required
// and should generate the AppCode field.
// In both cases the URLToken should be generated.
// The URLToken is directly used to create a clickable link in the email sent to user while mobile apps have to use the AppCode to exchange for the URLToken.
// SourceURL and AppCode are mutually exclusive.
type PwResetSession struct {
	// A long random string used to build a URL to be used on a web site.
	// It always exists even for mobile devices. To verify a session send from
	// mobiles apps, the request contains Email + AppCode to uniquely identify
	// a row in db since the AppCode is short and duplicate chances are very high.
	// Then the URLToken is sent back so that we could reset password
	// using URLToken, just like in a web page.
	Token      string      `json:"token" db:"token"`
	Email      string      `json:"email" db:"email"`
	SourceURL  null.String `json:"-" db:"source_url"` // Null for mobile apps
	IsUsed     bool        `json:"-" db:"is_used"`
	ExpiresIn  int64       `json:"-" db:"expires_in"`
	CreatedUTC chrono.Time `json:"-" db:"created_utc"`
}

// NewPwResetSession creates a new PwResetSession instance
// based on request body which contains a required `email`
// field, and an optionally `sourceUrl` field.
func NewPwResetSession(params input.ForgotPasswordParams) (PwResetSession, error) {
	token, err := gorest.RandomHex(32)
	if err != nil {
		return PwResetSession{}, err
	}

	if params.SourceURL.IsZero() {
		params.SourceURL = null.StringFrom("https://users.ftchinese.com/password-reset")
	}

	return PwResetSession{
		Email:      params.Email,
		SourceURL:  params.SourceURL,
		Token:      token, // URLToken always exists.
		IsUsed:     false,
		ExpiresIn:  3 * 60 * 60, // Valid for 3 hours
		CreatedUTC: chrono.TimeNow(),
	}, nil
}

// MustNewPwResetSession panic on error.
func MustNewPwResetSession(params input.ForgotPasswordParams) PwResetSession {
	s, err := NewPwResetSession(params)
	if err != nil {
		panic(err)
	}

	return s
}

// BuildURL creates password reset link.
// Returns an empty string if AppCode field exists so that
// the template will not render the URL section.
func (s PwResetSession) BuildURL() string {
	return fmt.Sprintf("%s/%s", s.SourceURL.String, s.Token)
}

// IsExpired tests whether an existing PwResetSession is expired.
func (s PwResetSession) IsExpired() bool {
	return s.CreatedUTC.Add(time.Second * time.Duration(s.ExpiresIn)).Before(time.Now())
}

func (s PwResetSession) DurHours() int64 {
	h := (time.Duration(s.ExpiresIn) * time.Second).Hours()
	return int64(h)
}

func (s PwResetSession) DurMinutes() int64 {
	m := (time.Duration(s.ExpiresIn) * time.Second).Minutes()
	return int64(m)
}

func (s PwResetSession) FormatDuration() string {
	if s.ExpiresIn >= 60*60 {
		return fmt.Sprintf("%d小时", s.DurHours())
	}

	return fmt.Sprintf("%d分钟", s.DurMinutes())
}
