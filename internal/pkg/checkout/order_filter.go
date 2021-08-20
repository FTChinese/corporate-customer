package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"strings"
)

// OrderFilter is used to hold the WHERE clause when
// retrieving a list orders.
// The WHERE clause differs based on how you use it:
// For the b2b app, an admin could only see order created by
// its team;
// For CMS, we need to see all orders created by any team.
type OrderFilter struct {
	TeamID string
	Status Status
}

// NewOrderFilter creates a default OrderFilter.
func NewOrderFilter(teamID string) OrderFilter {
	return OrderFilter{
		TeamID: teamID,
		Status: StatusNull,
	}
}

func (f OrderFilter) SQLWhere() pkg.SQLWhere {
	var clause strings.Builder
	var args = make([]interface{}, 0)

	if f.TeamID != "" {
		clause.WriteString("o.team_id = ?")
		args = append(args, f.TeamID)
	}

	if f.Status != StatusNull {
		clause.WriteString(" AND o.current_status = ?")
		args = append(args, f.Status.String())
	}

	return pkg.SQLWhere{
		Clause: clause.String(),
		Values: args,
	}
}
