package licence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type LicJSON struct {
	Licence
}

// Value implements Valuer interface by serializing an Invitation into
// JSON data.
func (l LicJSON) Value() (driver.Value, error) {
	if l.ID == "" {
		return nil, nil
	}

	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// Scan implements Valuer interface by deserializing an invitation field.
func (l *LicJSON) Scan(src interface{}) error {
	if src == nil {
		*l = LicJSON{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp LicJSON
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*l = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to LicJSON")
	}
}

// InvitationJSON is used to implement sql Valuer interface.
// Problems if you implement it on Invitation: when used
// as a field, the sql driver could save/retrieve a column
// as JSON; however, when you want to use Invitation as a
// plain SQL row, it continues to you custom `scan`, expecting
// JSON value instead of plain SQL columns.
type InvitationJSON struct {
	Invitation
}

// Value implements Valuer interface by serializing an Invitation into
// JSON data.
func (i InvitationJSON) Value() (driver.Value, error) {
	if i.ID == "" {
		return nil, nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// Scan implements Valuer interface by deserializing an invitation field.
func (i *InvitationJSON) Scan(src interface{}) error {
	if src == nil {
		*i = InvitationJSON{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp InvitationJSON
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*i = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to InvitationJSON")
	}
}
