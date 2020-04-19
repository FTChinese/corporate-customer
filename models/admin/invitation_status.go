package admin

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// InvitationStatus is an enumeration of invitation current phase.
type InvitationStatus int

const (
	InvitationStatusNull InvitationStatus = iota
	InvitationStatusCreated
	InvitationStatusAccepted
	InvitationStatusRevoked
)

var _invitationStatusNames = [...]string{
	"",
	"created",
	"accepted",
	"revoked",
}

// String representation of OrderKind
var _invitationStatusMap = map[InvitationStatus]string{
	1: _invitationStatusNames[1],
	2: _invitationStatusNames[2],
	3: _invitationStatusNames[3],
}

// Used to get OrderKind from a string.
var _invitationStatusValue = map[string]InvitationStatus{
	_invitationStatusNames[1]: 1,
	_invitationStatusNames[2]: 2,
	_invitationStatusNames[3]: 3,
}

// ParseInvitationStatus creates OrderKind from a string.
func ParseInvitationStatus(name string) (InvitationStatus, error) {
	if x, ok := _invitationStatusValue[name]; ok {
		return x, nil
	}

	return InvitationStatusNull, fmt.Errorf("%s is not valid InvitationStatus", name)
}

func (x InvitationStatus) String() string {
	if s, ok := _invitationStatusMap[x]; ok {
		return s
	}

	return ""
}

func (x *InvitationStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, _ := ParseInvitationStatus(s)

	*x = tmp

	return nil
}

func (x InvitationStatus) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

func (x *InvitationStatus) Scan(src interface{}) error {
	if src == nil {
		*x = InvitationStatusNull
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp, _ := ParseInvitationStatus(string(s))
		*x = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

func (x InvitationStatus) Value() (driver.Value, error) {
	s := x.String()
	if s == "" {
		return nil, nil
	}

	return s, nil
}
