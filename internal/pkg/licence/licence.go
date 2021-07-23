package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/dt"
	"github.com/FTChinese/ftacademy/pkg/price"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"time"
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
	ID string `json:"id" db:"licence_id"`
	price.Edition
	Creator
	CurrentPeriodStartUTC chrono.Time `json:"currentPeriodStartUtc" db:"current_period_start_utc"`
	CurrentPeriodEndUTC   chrono.Time `json:"currentPeriodEndUtc" db:"current_period_end_utc"`
	StartDateUTC          chrono.Time `json:"startDateUtc" db:"start_date_utc"`
	TrialStartUTC         chrono.Date `json:"trialStartUtc" db:"trial_start_utc"`
	TrialEndUTC           chrono.Date `json:"trialEndUtc" db:"trial_end_utc"`
	LatestOrderID         null.String `json:"latestOrderId" db:"latest_order_id"`
	LatestPrice           price.Price `json:"latestPrice" db:"latest_price"`

	// The following fields are subject to change.
	Status           Status      `json:"status" db:"lic_status"`
	LatestInvitation Invitation  `json:"latestInvitation" db:"latest_invitation"` // A redundant field that should be synced with the invitation row. I created this field since I don't know how to retrieve both licence and invitation row in SQL's flat structure, specially used when retrieve a list of such rows.
	AssigneeID       null.String `json:"-" db:"assignee_id"`
	RowMeta
}

func NewBaseLicence(p price.Price, orderID string, creator admin.PassportClaims) BaseLicence {

	now := time.Now()

	period := dt.NewTimeRange(now).
		WithCycle(p.Cycle)

	return BaseLicence{
		ID:      pkg.LicenceID(),
		Edition: p.Edition,
		Creator: Creator{
			CreatorID: creator.AdminID,
			TeamID:    creator.TeamID.String,
		},
		CurrentPeriodStartUTC: chrono.TimeFrom(period.Start),
		CurrentPeriodEndUTC:   chrono.TimeFrom(period.End),
		StartDateUTC:          chrono.TimeFrom(now),
		TrialStartUTC:         chrono.Date{},
		TrialEndUTC:           chrono.Date{},
		LatestOrderID:         null.StringFrom(orderID),
		LatestPrice:           p,
		Status:                LicStatusAvailable,
		LatestInvitation:      Invitation{},
		AssigneeID:            null.String{},
		RowMeta: RowMeta{
			CreatedUTC: chrono.TimeFrom(now),
			UpdatedUTC: chrono.TimeFrom(now),
		},
	}
}

// IsAvailable checks whether the licence is available to
// be assigned to a reader.
// As long as its status is not granted, it is available to
// be invited or assigned.
// An available licence should always be allowed to
// generate multiple invitations.
// However, only the latest invitation could be accepted.
// All previous invitations will be invalidated in such case.
func (l BaseLicence) IsAvailable() bool {
	return l.Status != LicStatusGranted && l.AssigneeID.IsZero()
}

// WithInvitation syncs a licence's latest_invitation when a new invitation is create for it.
func (l BaseLicence) WithInvitation(inv Invitation) BaseLicence {
	l.Status = LicStatusInvited
	l.LatestInvitation = inv
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// IsInvitationRevocable ensures that the licence could
// have its invitation revoked.
// A licence could only have its invitation revoked when
// an invitation is sent but not accepted.
func (l BaseLicence) IsInvitationRevocable() bool {
	return l.Status == LicStatusInvited && l.AssigneeID.IsZero()
}

// WithInvitationRevoked syncs the licence's invitation when
// the related invitation is revoked
// so that admin could invite another one to use this licence.
// This is used after invitation is sent but before it is
// accepted.
func (l BaseLicence) WithInvitationRevoked() BaseLicence {
	l.Status = LicStatusAvailable
	l.LatestInvitation = Invitation{}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

// Granted links a licence to a ftc reader.
func (l BaseLicence) Granted(a Assignee, inv Invitation) BaseLicence {
	l.Status = LicStatusGranted
	l.LatestInvitation = inv
	l.AssigneeID = a.FtcID
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

// IsGrantedTo checks if a licence is granted to the
// specified membership.
func (l BaseLicence) IsGrantedTo(m reader.Membership) bool {
	if !m.IsB2B() {
		return false
	}

	return m.B2BLicenceID.String == l.ID
}

func (l BaseLicence) IsRevocable() bool {
	return l.Status == LicStatusGranted && l.AssigneeID.Valid
}

// Revoked unlink a user from a licence.
func (l BaseLicence) Revoked() BaseLicence {
	l.Status = LicStatusAvailable
	l.LatestInvitation = Invitation{}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

// Licence contains data of a licence row joined with user info.
// This is used mostly when retrieving a list of licence.
type Licence struct {
	BaseLicence
	// Only exists after reader accepted an invitation.
	// Join with userinfo table
	Assignee Assignee `json:"assignee" db:"assignee"`
}

type LicList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []Licence `json:"data"`
	Err  error     `json:"-"`
}
