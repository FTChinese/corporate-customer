package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type Licence struct {
	ID         string      `db:"licence_id"`
	Tier       enum.Tier   `db:"tier"`
	Cycle      enum.Cycle  `db:"cycle"`
	ExpireDate chrono.Date `db:"expire_date"`
	AssigneeID null.String `db:"assignee_id"`
	CreatedUTC chrono.Time `db:"created_utc"`
	UpdatedUTC chrono.Time `db:"updated_utc"`
	Team       string      `db:"team"`
}
