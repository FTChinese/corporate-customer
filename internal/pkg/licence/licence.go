package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/dt"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"time"
)

// Licence represents a row in the licence table.
// A team can purchase multiple licences, as they usually do.
// After purchasing licences, the team can send
// invitations to team member for a specific licence.
// Once a team member accepts the invitation, this licence
// is assigned to this member and cannot be reassigned
// unless it is revoked.
// When a licence is revoked, backup reader's current
// membership status and delete that row, clear
// the AssigneeID field.
type Licence struct {
	ID string `json:"id" db:"licence_id"`
	price.Edition
	admin.Creator
	// The following fields are subject to change.
	Status                Status         `json:"status" db:"lic_status"`
	CurrentPeriodStartUTC chrono.Time    `json:"currentPeriodStartUtc" db:"current_period_start_utc"`
	CurrentPeriodEndUTC   chrono.Time    `json:"currentPeriodEndUtc" db:"current_period_end_utc"`
	StartDateUTC          chrono.Time    `json:"startDateUtc" db:"start_date_utc"`
	TrialStartUTC         chrono.Date    `json:"trialStartUtc" db:"trial_start_utc"`
	TrialEndUTC           chrono.Date    `json:"trialEndUtc" db:"trial_end_utc"`
	LatestOrderID         null.String    `json:"latestOrderId" db:"latest_order_id"`
	LatestPrice           price.Price    `json:"latestPrice" db:"latest_price"`
	LatestInvitation      InvitationJSON `json:"latestInvitation" db:"latest_invitation"` // A redundant field that should be synced with the invitation row. I created this field since I don't know how to retrieve both licence and invitation row in SQL's flat structure, specially used when retrieve a list of such rows.
	AssigneeID            null.String    `json:"-" db:"assignee_id"`
	admin.RowTime
}

func NewLicence(p price.Price, orderID string, creator admin.Creator) Licence {

	// TODO: pass as parameter
	now := time.Now()

	period := dt.NewTimeRange(now).
		WithCycle(p.Cycle)

	return Licence{
		ID:                    pkg.LicenceID(),
		Edition:               p.Edition,
		Creator:               creator,
		CurrentPeriodStartUTC: chrono.TimeUTCFrom(period.Start),
		CurrentPeriodEndUTC:   chrono.TimeUTCFrom(period.End),
		StartDateUTC:          chrono.TimeUTCFrom(now),
		TrialStartUTC:         chrono.Date{},
		TrialEndUTC:           chrono.Date{},
		LatestOrderID:         null.StringFrom(orderID),
		LatestPrice:           p,
		Status:                LicStatusAvailable,
		LatestInvitation:      InvitationJSON{},
		AssigneeID:            null.String{},
		RowTime:               admin.NewRowTime(),
	}
}

func (l Licence) RenewalStartTime() time.Time {
	now := time.Now()

	if l.CurrentPeriodEndUTC.After(now) {
		return l.CurrentPeriodEndUTC.Time
	}

	return now
}

// IsAvailable checks whether the licence is available to
// be assigned to a reader.
// As long as its status is not granted, it is available to
// be invited or assigned.
// An available licence should always be allowed to
// generate multiple invitations.
// However, only the latest invitation could be accepted.
// All previous invitations will be invalidated in such case.
func (l Licence) IsAvailable() bool {
	return l.Status != LicStatusGranted && l.AssigneeID.IsZero()
}

// WithInvitation syncs a licence's latest_invitation when a new invitation is create for it.
func (l Licence) WithInvitation(inv Invitation) Licence {
	l.Status = LicStatusInvited
	l.LatestInvitation = InvitationJSON{inv}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// IsInvitationRevocable ensures that the licence could
// have its invitation revoked.
// A licence could only have its invitation revoked when
// an invitation is sent but not accepted.
func (l Licence) IsInvitationRevocable() bool {
	return l.Status == LicStatusInvited && l.AssigneeID.IsZero()
}

// WithInvitationRevoked syncs the licence's invitation when
// the related invitation is revoked
// so that admin could invite another one to use this licence.
// This is used after invitation is sent but before it is
// accepted.
func (l Licence) WithInvitationRevoked() Licence {
	l.Status = LicStatusAvailable
	l.LatestInvitation = InvitationJSON{}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeUTCNow()

	return l
}

// Granted links a licence to a ftc reader.
func (l Licence) Granted(a Assignee, inv Invitation) Licence {
	l.Status = LicStatusGranted
	l.LatestInvitation = InvitationJSON{inv}
	l.AssigneeID = a.FtcID
	l.UpdatedUTC = chrono.TimeUTCNow()

	return l
}

// IsGranted tests whether a licence is granted to any reader.
func (l Licence) IsGranted() bool {
	return l.Status == LicStatusGranted && l.AssigneeID.Valid
}

// IsGrantedTo checks if a licence is granted to the
// specified membership.
func (l Licence) IsGrantedTo(m reader.Membership) bool {
	if !m.IsB2B() {
		return false
	}

	return l.ID == m.B2BLicenceID.String
}

func (l Licence) IsRevocable() bool {
	return l.Status == LicStatusGranted && l.AssigneeID.Valid
}

// Revoked unlink a user from a licence.
func (l Licence) Revoked() Licence {
	l.Status = LicStatusAvailable
	l.LatestInvitation = InvitationJSON{}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeUTCNow()

	return l
}

func (l Licence) Renewed(p price.Price, orderID string) Licence {
	now := time.Now()

	startTime := l.RenewalStartTime()
	period := dt.NewTimeRange(startTime).WithCycle(p.Cycle)

	l.CurrentPeriodStartUTC = chrono.TimeUTCFrom(period.Start)
	l.CurrentPeriodEndUTC = chrono.TimeUTCFrom(period.End)
	l.LatestPrice = p
	l.LatestOrderID = null.StringFrom(orderID)
	l.UpdatedUTC = chrono.TimeUTCFrom(now)

	return l
}

// ExpandedLicence contains data of a licence row joined with user info.
// This is used mostly when retrieving a list of licence.
type ExpandedLicence struct {
	Licence
	// Only exists after reader accepted an invitation.
	// Join with userinfo table
	Assignee AssigneeJSON `json:"assignee" db:"assignee"`
}

type PagedLicenceList struct {
	pkg.PagedList
	Data []ExpandedLicence `json:"data"`
}
