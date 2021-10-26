package admin

import (
	"github.com/FTChinese/ftacademy/pkg/oauth"
	"log"
	"net/http"
)

type JWTGuard struct {
	signingKey []byte
}

func NewJWTGuard(key []byte) JWTGuard {
	return JWTGuard{
		signingKey: key,
	}
}

func (g JWTGuard) GetKey() []byte {
	return g.signingKey
}

func (g JWTGuard) CreatePassport(a BaseAccount) (Passport, error) {
	return NewPassport(a, g.signingKey)
}

func (g JWTGuard) RetrievePassportClaims(req *http.Request) (PassportClaims, error) {
	ss, err := oauth.GetBearerAuth(req.Header)
	if err != nil {
		log.Printf("Error parsing Authorization header: %v", err)
		return PassportClaims{}, err
	}

	claims, err := ParsePassportClaims(ss, g.signingKey)
	if err != nil {
		log.Printf("Error parsing JWT %v", err)
		return PassportClaims{}, err
	}

	return claims, nil
}
