package admin

import "github.com/guregu/null"

// Assignee represents a reader who can accept
// an invitation email, and who can be granted
// a licence.
type Assignee struct {
	Email    null.String `db:"email"`
	FtcID    null.String `db:"ftc_id"`
	UserName null.String `db:"user_name"`
	IsVIP    bool        `db:"is_vip"`
}

type AssigneeSchema struct {
	Assignee
	TeamID string `db:"team_id"`
}
