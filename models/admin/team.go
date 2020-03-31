package admin

import (
	"github.com/FTChinese/b2b/models/validator"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"strings"
)

// TeamForm is the form data submitted to create a Team.
type TeamForm struct {
	Name         string `form:"name"`         // Organization name. This is used for display only. Required. 128 chars.
	InvoiceTitle string `form:"invoiceTitle"` // Invoice title might be different from Name. Optional.
	Errors       map[string]string
}

func (f *TeamForm) Validate() bool {
	f.Name = strings.TrimSpace(f.Name)
	f.InvoiceTitle = strings.TrimSpace(f.InvoiceTitle)

	msg := validator.New("机构名称").Required().Max(128).Validate(f.Name)
	if msg != "" {
		f.Errors["name"] = msg
	}

	msg = validator.New("发票抬头").Max(512).Validate(f.InvoiceTitle)
	if msg != "" {
		f.Errors["invoiceTitle"] = msg
	}

	return len(f.Errors) == 0
}

// Team represents an existing b2b entity.
// An admin account can create teams.
// A team can purchase licences.
type Team struct {
	ID           string      `db:"team_id"`
	Name         string      `db:"name"`
	InvoiceTitle null.String `db:"invoice_title"`
	CreatedUTC   chrono.Time `db:"created_utc"`
	UpdatedUTC   chrono.Time `db:"updated_utc"`
	AdminID      string      `db:"admin_id"`
}

func NewTeam(f TeamForm, adminID string) Team {
	return Team{
		ID:           "team_" + rand.String(12),
		Name:         f.Name,
		InvoiceTitle: null.NewString(f.InvoiceTitle, f.InvoiceTitle != ""),
		CreatedUTC:   chrono.TimeNow(),
		UpdatedUTC:   chrono.TimeNow(),
		AdminID:      adminID,
	}
}

// IsEqual compares the original values of team and
// updated values. If they are equal, no further operation
// should be conducted.
func (t Team) IsEqual(f TeamForm) bool {
	return t.Name == f.Name && t.InvoiceTitle.String == f.InvoiceTitle
}

// Update updates the current team to new values.
func (t *Team) Update(f TeamForm) {
	t.Name = f.Name
	t.InvoiceTitle = null.NewString(f.InvoiceTitle, f.InvoiceTitle != "")
}
