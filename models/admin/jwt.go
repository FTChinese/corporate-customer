package admin

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/guregu/null"
	"log"
	"time"
)

// AccountClaims is a JWT custom claims containing only the
// essential fields of an account so that the signed string
// won't become too long while the backend can determine
// user's identity.
// After user logged in, the JWT is send to client as one
// of the JSON fields. The response body contains more fields
// than this claims so that client is able to show extra
// information on UI.
type AccountClaims struct {
	AdminID string      `json:"aid"`
	TeamID  null.String `json:"tid"`
	jwt.StandardClaims
}

// SignedString create a JWT based on current claims.
func (c AccountClaims) SignedString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	ss, err := token.SignedString(signingKey)

	if err != nil {
		return ss, err
	}

	return ss, nil
}

// JWTAccount adds ExpiresAt so that client
// could check whether the login session is
// expired. It carries the Json Web Token.
type JWTAccount struct {
	Account
	TeamID    null.String `json:"teamId" db:"team_id"`
	ExpiresAt int64       `json:"expiresAt"`
	Token     string      `json:"token"`
}

func (a JWTAccount) Claims() AccountClaims {
	return AccountClaims{
		AdminID: a.ID,
		TeamID:  a.TeamID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400*7,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "com.ftchinese.b2b",
		},
	}
}

func (a JWTAccount) WithToken() (JWTAccount, error) {
	claims := a.Claims()

	ss, err := claims.SignedString()
	if err != nil {
		return JWTAccount{}, err
	}

	a.ExpiresAt = claims.ExpiresAt
	a.Token = ss
	return a, nil
}

func ParseJWT(ss string) (AccountClaims, error) {
	token, err := jwt.ParseWithClaims(
		ss,
		&AccountClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return signingKey, nil
		})

	if err != nil {
		log.Printf("Parsing JWT error: %v", err)
		return AccountClaims{}, err
	}

	log.Printf("Claims: %v", token.Claims)

	// NOTE: token.Claims is an interface, so it is a pointer, not a value type.
	if claims, ok := token.Claims.(*AccountClaims); ok {
		return *claims, nil
	}
	return AccountClaims{}, errors.New("wrong JWT claims")
}
