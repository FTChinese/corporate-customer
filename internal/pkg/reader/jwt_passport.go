package reader

import (
	"errors"
	"github.com/FTChinese/go-rest/enum"
	"github.com/golang-jwt/jwt"
	"github.com/guregu/null"
	"log"
	"time"
)

func NewStandardClaims(expiresAt int64) jwt.StandardClaims {
	return jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "cn.ftacademy.reader",
	}
}

// passportVersion determines whether client is using the latest
// passport claims data structure. Invalidate a session if not.
const passportVersion = 2

type PassportClaims struct {
	FtcID       string           `json:"fid"`
	UnionID     null.String      `json:"wid"`
	LoginMethod enum.LoginMethod `json:"mtd"`
	Live        bool             `json:"live"` // Whether the account should use live API.
	Version     int              `json:"v"`
	jwt.StandardClaims
}

func (c PassportClaims) VersionMatched() bool {
	return c.Version == passportVersion
}

type Passport struct {
	Account
	ExpiresAt int64  `json:"expiresAt"`
	Live      bool   `json:"live"`
	Version   int    `json:"version"`
	Token     string `json:"token"`
}

func NewPassport(a Account, signingKey []byte) (Passport, error) {
	claims := PassportClaims{
		FtcID:          a.FtcID,
		UnionID:        a.UnionID,
		LoginMethod:    a.LoginMethod,
		Live:           !a.IsTest(),
		Version:        passportVersion,
		StandardClaims: NewStandardClaims(time.Now().Unix() + 86400*7),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(signingKey)

	if err != nil {
		return Passport{}, err
	}

	return Passport{
		Account:   a,
		ExpiresAt: claims.ExpiresAt,
		Live:      claims.Live,
		Version:   claims.Version,
		Token:     ss,
	}, nil
}

func ParsePassportClaims(ss string, key []byte) (PassportClaims, error) {
	token, err := jwt.ParseWithClaims(
		ss,
		&PassportClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return key, nil
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
