package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
)

type TeamJSON struct {
	input.TeamParams
}

func (t TeamJSON) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (t *TeamJSON) Scan(src interface{}) error {
	if src == nil {
		*t = TeamJSON{}
	}
	switch s := src.(type) {
	case []byte:
		var tmp TeamJSON
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}

		*t = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to TeamJSON")
	}
}

// CMSOrderRow is a row in the order table that is slightly
// different from the version used by admin.
type CMSOrderRow struct {
	OrderRow
	Team TeamJSON `json:"team" db:"team"`
}

type CMSOrderList struct {
	pkg.PagedList
	Data []CMSOrderRow `json:"data"`
}
