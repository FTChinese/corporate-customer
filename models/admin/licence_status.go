package admin

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type LicenceStatus int

const (
	LicStatusNull LicenceStatus = iota
	LicStatusAvailable
	LicStatusInvited
	LicStatusGranted
)

var _licenceStatusNames = [...]string{
	"",
	"available",
	"invited",
	"granted",
}

// String representation of OrderKind
var _licenceStatusMap = map[LicenceStatus]string{
	1: _licenceStatusNames[1],
	2: _licenceStatusNames[2],
	3: _licenceStatusNames[3],
}

// Used to get OrderKind from a string.
var _licenceStatusValue = map[string]LicenceStatus{
	_licenceStatusNames[1]: 1,
	_licenceStatusNames[2]: 2,
	_licenceStatusNames[3]: 3,
}

// ParseLicenceStatus creates OrderKind from a string.
func ParseLicenceStatus(name string) (LicenceStatus, error) {
	if x, ok := _licenceStatusValue[name]; ok {
		return x, nil
	}

	return LicStatusNull, fmt.Errorf("%s is not valid LicenceStatus", name)
}

func (x LicenceStatus) String() string {
	if s, ok := _licenceStatusMap[x]; ok {
		return s
	}

	return ""
}

func (x *LicenceStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, _ := ParseLicenceStatus(s)

	*x = tmp

	return nil
}

func (x LicenceStatus) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

func (x *LicenceStatus) Scan(src interface{}) error {
	if src == nil {
		*x = LicStatusNull
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp, _ := ParseLicenceStatus(string(s))
		*x = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

func (x LicenceStatus) Value() (driver.Value, error) {
	s := x.String()
	if s == "" {
		return nil, nil
	}

	return s, nil
}
