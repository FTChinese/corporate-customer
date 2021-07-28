package licence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/guregu/null"
	"strings"
)

// Assignee represents a reader who can accept
// an invitation email, and who can be granted
// a licence.
type Assignee struct {
	FtcID    null.String `json:"id" db:"ftc_id"`
	UnionID  null.String `json:"unionId" db:"wx_union_id"`
	Email    null.String `json:"email" db:"user_email"`
	UserName null.String `json:"userName" db:"user_name"`
}

func (a Assignee) IsZero() bool {
	return a.FtcID.IsZero()
}

// NormalizeName tries to find a proper way to greet user
// in email.
func (a Assignee) NormalizeName() string {
	if a.UserName.Valid {
		return a.UserName.String
	}

	return strings.Split(a.Email.String, "@")[0]
}

func (a Assignee) TeamMember(teamID string) Staffer {
	return Staffer{
		Email:  a.Email.String,
		TeamID: teamID,
		FtcID:  a.FtcID,
	}
}

// AssigneeJSON is used to save Assignee as JSON value.
type AssigneeJSON struct {
	Assignee
}

func (a AssigneeJSON) Value() (driver.Value, error) {
	if a.FtcID.IsZero() {
		return nil, nil
	}

	return json.Marshal(a)
}

func (a *AssigneeJSON) Scan(src interface{}) error {
	if src == nil {
		*a = AssigneeJSON{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp AssigneeJSON
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*a = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}
