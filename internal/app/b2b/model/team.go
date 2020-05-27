package model

import (
	"github.com/FTChinese/ftacademy/internal/pkg/validator"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"strings"
)

// Team represents an existing b2b entity.
// An admin account can create teams.
// A team can purchase licences.
type Team struct {
	ID           string      `db:"team_id"`
	AdminID      string      `db:"team_id"`
	Name         string      `json:"name" db:"name"`
	InvoiceTitle null.String `json:"invoiceTitle" db:"invoice_title"`
	CreatedUTC   chrono.Time `db:"created_utc"`
	UpdatedUTC   chrono.Time `db:"updated_utc"`
}

func (t *Team) Validate() *render.ValidationError {
	t.Name = strings.TrimSpace(t.Name)
	title := strings.TrimSpace(t.InvoiceTitle.String)
	t.InvoiceTitle = null.NewString(title, title != "")

	ve := validator.New("name").Required().Max(128).Validate(t.Name)
	if ve != nil {
		return ve
	}

	return validator.New("invoiceTitle").Max(512).Validate(title)
}

func (t Team) BuildOn(adminID string) Team {
	t.ID = "team_" + rand.String(12)
	t.AdminID = adminID
	t.CreatedUTC = chrono.TimeNow()
	t.UpdatedUTC = chrono.TimeNow()
	return t
}

// IsEqual compares the original values of team and
// updated values. If they are equal, no further operation
// should be conducted.
func (t Team) IsEqual(newVal Team) bool {
	return t.Name == newVal.Name && t.InvoiceTitle == newVal.InvoiceTitle
}

// Update updates the current team to new values.
func (t Team) Update(newVal Team) Team {
	t.Name = newVal.Name
	t.InvoiceTitle = newVal.InvoiceTitle
	t.UpdatedUTC = chrono.TimeNow()

	return t
}
