package repository

import "errors"

var (
	ErrLicenceUnavailable = errors.New("licence is unavailable to grant")
	ErrAlreadyMember      = errors.New("the invitee already has a valid membership")
	ErrInviteeMismatch    = errors.New("an invitation for this licence is already sent to another user")
	ErrInvalidInvitation  = errors.New("invalid invitation")
	ErrLicenceTaken       = errors.New("the licence is already taken by another user")
)
