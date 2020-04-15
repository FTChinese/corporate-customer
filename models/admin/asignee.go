package admin

import (
	"github.com/guregu/null"
	"strings"
)

// Assignee represents a reader who can accept
// an invitation email, and who can be granted
// a licence.
type Assignee struct {
	FtcID    null.String `json:"ftcId" db:"ftc_id"`
	Email    null.String `json:"email" db:"email"`
	UserName null.String `json:"userName" db:"user_name"`
	IsVIP    bool        `json:"isVip" db:"is_vip"`
}

// NormalizeName tries to find a proper way to greet user
// in email.
func (a Assignee) NormalizeName() string {
	if a.UserName.Valid {
		return a.UserName.String
	}

	return strings.Split(a.Email.String, "@")[0]
}

type AssigneeSchema struct {
	Assignee
	TeamID string `db:"team_id"`
}
