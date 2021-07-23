package licence

import (
	"encoding/json"
	"errors"
	"github.com/guregu/null"
	"strings"
)

// Assignee represents a reader who can accept
// an invitation email, and who can be granted
// a licence.
type Assignee struct {
	FtcID    null.String `json:"ftcId" db:"ftc_id"`
	UnionID  null.String `json:"unionId" db:"wx_union_id"`
	Email    null.String `json:"email" db:"user_email"`
	UserName null.String `json:"userName" db:"user_name"`
}

func (a Assignee) IsZero() bool {
	return a.FtcID.IsZero()
}

//func (a Assignee) Value() (driver.Value, error) {
//	if a.FtcID.IsZero() {
//		return nil, nil
//	}
//
//	return a.FtcID.String, nil
//}

func (a *Assignee) Scan(src interface{}) error {
	if src == nil {
		*a = Assignee{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp Assignee
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
