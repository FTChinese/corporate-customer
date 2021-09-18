package subs

import "errors"

var (
	ErrOverrideAutoRenewForbidden = errors.New("cannot override a valid auto-renewal subscription")
)
