package admin

import (
	"database/sql/driver"
	"errors"
)

type LicenceStatus string

const (
	LicStatusNull      LicenceStatus = ""
	LicStatusAvailable LicenceStatus = "available"
	LicStatusInvited   LicenceStatus = "invitation_sent"
	LicStatusGranted   LicenceStatus = "granted"
)

func (x LicenceStatus) MarshalJSON() ([]byte, error) {

	if x == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + x + `"`), nil
}

func (x *LicenceStatus) Scan(src interface{}) error {
	if src == nil {
		*x = LicStatusNull
		return nil
	}

	switch s := src.(type) {
	case []byte:
		*x = LicenceStatus(s)
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

func (x LicenceStatus) Value() (driver.Value, error) {
	if x == "" {
		return nil, nil
	}

	return x, nil
}
