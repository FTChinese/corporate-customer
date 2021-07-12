package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
)

// Team represents an existing b2b entity.
// An admin account can create teams.
// A team can purchase licences.
type Team struct {
	ID      string `db:"team_id"`
	AdminID string `db:"team_id"`
	input.TeamParams
	CreatedUTC chrono.Time `db:"created_utc"`
	UpdatedUTC chrono.Time `db:"updated_utc"`
}

func NewTeam(adminID string, params input.TeamParams) Team {
	return Team{
		ID:         "team_" + rand.String(12),
		AdminID:    adminID,
		TeamParams: params,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.Time{},
	}
}

// Update updates the current team to new values.
func (t Team) Update(newVal input.TeamParams) Team {
	t.OrgName = newVal.OrgName
	t.InvoiceTitle = newVal.InvoiceTitle
	t.UpdatedUTC = chrono.TimeNow()

	return t
}
