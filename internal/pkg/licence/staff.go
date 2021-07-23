package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/guregu/null"
)

// Staffer is a member belong to a team under admin's
// management.
type Staffer struct {
	ID     int64       `json:"id" db:"id"`
	Email  string      `json:"email" db:"email"`
	TeamID string      `json:"teamId" db:"team_id"`
	FtcID  null.String `json:"ftcId" db:"ftc_id"`
}

// StaffList contains a list of assignee rows and the total number of rows for current team.
type StaffList struct {
	pkg.PagedList
	Data []Staffer `json:"data"`
}
