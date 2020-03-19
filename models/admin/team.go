package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"strings"
)

type TeamForm struct {
	Name         string `form:"name"`
	InvoiceTitle string `form:"invoiceTitle"`
	Errors       map[string]string
}

func (f *TeamForm) Sanitize() *TeamForm {
	f.Name = strings.TrimSpace(f.Name)
	f.InvoiceTitle = strings.TrimSpace(f.InvoiceTitle)

	return f
}

func (f *TeamForm) Validate() bool {
	return len(f.Errors) == 0
}

func (f *TeamForm) BuildTeam() Team {
	return Team{
		ID:           "",
		Name:         f.Name,
		InvoiceTitle: null.NewString(f.InvoiceTitle, f.InvoiceTitle != ""),
		CreatedUTC:   chrono.TimeNow(),
		UpdatedUTC:   chrono.TimeNow(),
	}
}

type Team struct {
	ID           string      `db:"team_id"`
	Name         string      `db:"name"`
	InvoiceTitle null.String `db:"invoice_title"`
	CreatedUTC   chrono.Time `db:"created_utc"`
	UpdatedUTC   chrono.Time `db:"updated_utc"`
	Admin        string      `db:"admin"`
}
