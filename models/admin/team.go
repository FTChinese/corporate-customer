package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type Team struct {
	ID           string      `db:"team_id"`
	Name         string      `db:"name"`
	InvoiceTitle null.String `db:"invoice_title"`
	CreatedUTC   chrono.Time `db:"created_utc"`
	UpdatedUTC   chrono.Time `db:"updated_utc"`
	Admin        string      `db:"admin"`
}
