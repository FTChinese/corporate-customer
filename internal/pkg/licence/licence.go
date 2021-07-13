package licence

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/pkg/price"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// BaseLicence is the licence a team purchased.
// A team can purchase multiple licences, as they usually do.
// After purchasing licences, the team can send
// invitations to team member for a specific licence.
// Once a team member accepts the invitation, this licence
// is assigned to this member and cannot be reassigned
// unless it is revoked.
// When a licence is revoked, backup reader's current
// current membership status and delete that row, clear
// the AssigneeID field.
type BaseLicence struct {
	ID     string `json:"id" db:"licence_id"`
	TeamID string `json:"teamId" db:"team_id"`
	price.Edition
	CreatedUTC            chrono.Time `json:"createdUtc" db:"created_utc"`
	Status                Status      `json:"status" db:"licence_status"`
	CurrentPeriodStartUTC chrono.Time `json:"currentPeriodStartUtc" db:"current_period_start_utc"`
	CurrentPeriodEndUTC   chrono.Time `json:"currentPeriodEndUtc" db:"current_period_end_utc"`
	StartDateUTC          chrono.Time `json:"startDateUtc" db:"start_date_utc"`
	TrialStartUTC         chrono.Date `json:"trialStartUtc" db:"trial_start_utc"`
	TrialEndUTC           chrono.Date `json:"trialEndUtc" db:"trial_end_utc"`
	LatestInvoiceID       null.String `json:"latestInvoiceId" db:"latest_invoice_id"`
	UpdatedUTC            chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

// Licence contains all data of a licence.
type Licence struct {
	BaseLicence
	// Only exists after reader accepted an invitation.
	// Join with userinfo table
	Assignee         Assignee    `json:"assignee" db:"assignee"`
	LatestPrice      price.Price `json:"latestPrice" db:"latest_price"`
	LatestInvitation Invitation  `json:"latestInvitation" db:"latest_invitation_id"`
}

// IsAvailable checks whether the licence is
// assigned to someone else.
func (l Licence) IsAvailable() bool {
	return l.Status == LicStatusAvailable && l.Assignee.FtcID.IsZero()
}

// WithInvited updates invitation status to invited after an invitation is sent.
func (l Licence) WithInvited(inv Invitation) Licence {
	l.Status = LicStatusInvited
	l.LatestInvitation = inv
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// CanInvitationBeRevoked ensures that the licence could
// have its invitation revoked.
// A licence could only have its invitation revoked when
// an invitation is sent but not accepted.
func (l Licence) CanInvitationBeRevoked() bool {
	return l.Status == LicStatusInvited && l.Assignee.FtcID.IsZero()
}

// WithInvitationRevoked revokes an invitation of a licence
// so that admin could invite another one to use this licence.
// This is used after invitation is sent but before it is
// accepted.
func (l Licence) WithInvitationRevoked() Licence {
	l.Status = LicStatusAvailable
	l.LatestInvitation = Invitation{}
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

// Revoke unlink a user from a licence.
func (l Licence) Revoke() Licence {
	l.Assignee = Assignee{}
	l.Status = LicStatusAvailable
	l.LatestInvitation = Invitation{}
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

func (l Licence) CanBeGranted() bool {
	return l.Status == LicStatusInvited && l.Assignee.FtcID.IsZero()
}

// WithGranted links a licence to a ftc reader.
func (l Licence) WithGranted(a Assignee) Licence {
	l.Assignee = a
	l.Status = LicStatusGranted
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// LicSchema defines the DB schema for licence table.
// Deprecated
type LicSchema struct {
	BaseLicence
	Assignee                     // Used to scan data from db.
	LatestPrice      null.String `db:"latest_price"` // The raw JSON column is a string.
	LatestInvitation null.String `db:"latest_invitation"`
}

// ToLicence converts the data retrieved from db to Licence instance.
// Deprecated
func (ls LicSchema) ToLicence() (Licence, error) {
	var p price.Price
	var inv Invitation

	if ls.LatestPrice.Valid {
		err := json.Unmarshal([]byte(ls.LatestPrice.String), &p)
		if err != nil {
			return Licence{}, err
		}
	}

	if ls.LatestInvitation.Valid {
		err := json.Unmarshal([]byte(ls.LatestInvitation.String), &inv)
		if err != nil {
			return Licence{}, err
		}
	}
	return Licence{
		BaseLicence:      ls.BaseLicence,
		Assignee:         ls.Assignee,
		LatestPrice:      p,
		LatestInvitation: inv,
	}, nil
}

type LicList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []Licence `json:"data"`
	Err  error     `json:"-"`
}

// ExpandedLicence includes the data of the licence,
// its attached plan, a latest invitation if present
// and the assignee if it is already granted.
// Deprecated
type ExpandedLicence struct {
	Licence
	Assignee Assignee `json:"assignee"` // If no use is granted to use this licence, its fields are empty.
}

// ExpLicenceSchema is used to save/retrieve licence.
// Deprecated
type ExpLicenceSchema struct {
	LicSchema
	Assignee
}

// ExpandedLicence transforms a row retrieved from DB to output format.
// Deprecated
func (s ExpLicenceSchema) ExpandedLicence() (ExpandedLicence, error) {

	lic, err := s.ToLicence()
	if err != nil {
		return ExpandedLicence{}, err
	}

	return ExpandedLicence{
		Licence:  lic,
		Assignee: s.Assignee,
	}, nil
}

// PagedExpLicences is used to output a list of licence with pagination.
// Deprecated
type PagedExpLicences struct {
	Total int64             `json:"total"` // The total number of rows.
	Data  []ExpandedLicence `json:"data"`
	Err   error             `json:"-"`
}
