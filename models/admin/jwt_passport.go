package admin

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/guregu/null"
	"log"
	"time"
)

func NewStandardClaims(expiresAt int64) jwt.StandardClaims {
	return jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "com.ftchinese.b2b",
	}
}

// PassportClaims is a JWT custom claims containing only the
// essential fields of an account so that the signed string
// won't become too long while the backend can determine
// user's identity.
// After user logged in, the JWT is send to client as one
// of the JSON fields. The response body contains more fields
// than this claims so that client is able to show extra
// information on UI.
type PassportClaims struct {
	AdminID string      `json:"aid"`
	TeamID  null.String `json:"tid"`
	jwt.StandardClaims
}

func NewPassportClaims(a Passport) PassportClaims {
	return PassportClaims{
		AdminID:        a.ID,
		TeamID:         a.TeamID,
		StandardClaims: NewStandardClaims(time.Now().Unix() + 86400*7),
	}
}

// SignedString create a JWT based on current claims.
func (c PassportClaims) SignedString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	ss, err := token.SignedString(signingKey)

	if err != nil {
		return ss, err
	}

	return ss, nil
}

type Passport struct {
	Account
	TeamID   null.String `json:"teamId" db:"team_id"`
	TeamName null.String `json:"-" db:"team_name"`
}

// PassportBearer carries the Json Web Token for a logged in
// admin plus structured data of it so that client do not
// need to decode the encoded data.
type PassportBearer struct {
	Passport
	ExpiresAt int64  `json:"expiresAt"`
	Token     string `json:"token"`
}

func NewPassportBearer(pp Passport) (PassportBearer, error) {
	claims := NewPassportClaims(pp)

	ss, err := claims.SignedString()
	if err != nil {
		return PassportBearer{}, err
	}

	return PassportBearer{
		Passport:  pp,
		ExpiresAt: claims.ExpiresAt,
		Token:     ss,
	}, nil
}

func ParsePassportClaims(ss string) (PassportClaims, error) {
	token, err := jwt.ParseWithClaims(
		ss,
		&PassportClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return signingKey, nil
		})

	if err != nil {
		log.Printf("Parsing JWT error: %v", err)
		return PassportClaims{}, err
	}

	log.Printf("Claims: %v", token.Claims)

	// NOTE: token.Claims is an interface, so it is a pointer, not a value type.
	if claims, ok := token.Claims.(*PassportClaims); ok {
		return *claims, nil
	}
	return PassportClaims{}, errors.New("wrong JWT claims")
}
