package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Status int

const (
	StatusNull Status = iota
	StatusPending
	StatusPaid
	StatusProcessing
	StatusCancelled
)

var statusNames = [...]string{
	"",
	"pending_payment",
	"paid",
	"processing",
	"cancelled",
}

var statusMap = map[Status]string{
	1: statusNames[1],
	2: statusNames[2],
	3: statusNames[3],
	4: statusNames[4],
}

var statusValue = map[string]Status{
	statusNames[1]: 1,
	statusNames[2]: 2,
	statusNames[3]: 3,
	statusNames[4]: 4,
}

func ParseStatus(name string) (Status, error) {
	if x, ok := statusValue[name]; ok {
		return x, nil
	}

	return StatusNull, fmt.Errorf("%s is not a valid Status", name)
}

func (x Status) String() string {
	if s, ok := statusMap[x]; ok {
		return s
	}

	return ""
}

func (x *Status) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, _ := ParseStatus(s)

	*x = tmp

	return nil
}

// MarshalJSON implements the Marshaler interface
func (x Status) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

// Scan implements sql.Scanner interface to retrieve value from SQL.
// SQL null will be turned into zero value TierFree.
func (x *Status) Scan(src interface{}) error {
	if src == nil {
		*x = StatusNull
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp, _ := ParseStatus(string(s))
		*x = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to Status")
	}
}

// Value implements driver.Valuer interface to save value into SQL.
func (x Status) Value() (driver.Value, error) {
	s := x.String()
	if s == "" {
		return nil, nil
	}

	return s, nil
}
