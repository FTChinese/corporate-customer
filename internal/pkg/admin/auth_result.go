package admin

// AuthResult contains the data retrieved from db after
// authenticating password for the specified email.
// If db returns ErrNoRows, it indicates the specified email does not
// exists.
// The the PasswordMatched field is false, it indicates the account
// exists but password is not correct.
type AuthResult struct {
	AdminID         string `db:"id"`
	PasswordMatched bool   `db:"password_matched"`
}
