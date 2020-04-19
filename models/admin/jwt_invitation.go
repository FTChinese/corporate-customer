package admin

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type InvitationBearer struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewInvitationBearer(claims InviteeClaims) (InvitationBearer, error) {
	ss, err := claims.SignedString()
	if err != nil {
		return InvitationBearer{}, err
	}

	return InvitationBearer{
		Email: claims.Email,
		Token: ss,
	}, nil
}

type InviteeClaims struct {
	InvitationID string `json:"invId"`
	LicenceID    string `json:"licId"`
	TeamID       string `json:"teamId"`
	Email        string `json:"email"`
	FtcID        string `json:"ftcId"`
	jwt.StandardClaims
}

func NewInviteeClaims(inv Invitation) InviteeClaims {
	return InviteeClaims{
		InvitationID:   inv.ID,
		LicenceID:      inv.LicenceID,
		TeamID:         inv.TeamID,
		Email:          inv.Email,
		FtcID:          "",
		StandardClaims: NewStandardClaims(time.Now().Unix() + 600),
	}
}

func (c InviteeClaims) SignedString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	ss, err := token.SignedString(signingKey)

	if err != nil {
		return ss, err
	}

	return ss, nil
}

func ParseInviteeClaims(ss string) (InviteeClaims, error) {
	token, err := jwt.ParseWithClaims(
		ss,
		&InviteeClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})

	if err != nil {
		return InviteeClaims{}, err
	}

	if claims, ok := token.Claims.(*InviteeClaims); ok {
		return *claims, nil
	}

	return InviteeClaims{}, errors.New("wrong JWT claims")
}
