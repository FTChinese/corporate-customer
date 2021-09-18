package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
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
	StartDateUTC          chrono.Time    `json:"startDateUtc" db:"start_date_utc"` // Initial start time of licence become effective for the first time. Might be different from CreatedUTC.
	TrialStartUTC         chrono.Date    `json:"trialStartUtc" db:"trial_start_utc"`
	TrialEndUTC           chrono.Date    `json:"trialEndUtc" db:"trial_end_utc"`
	HintGrantMismatch     bool           `json:"hintGrantMismatch" db:"hint_grant_mismatch"` // Indicates possible mismatch between the licence and the granted membership. It could happen upon renewal, or in a periodic verification.
	LatestTransactionID   null.String    `json:"latestTransactionId" db:"latest_transaction_id"`
	LatestPrice           price.Price    `json:"latestPrice" db:"latest_price"`
	LatestInvitation      InvitationJSON `json:"latestInvitation" db:"latest_invitation"` // A redundant field that should be synced with the invitation row. I created this field since I don't know how to retrieve both licence and invitation row in SQL's flat structure, specially used when retrieve a list of such rows.
	AssigneeID            null.String    `json:"assigneeId" db:"assignee_id"`
	admin.RowTime
}

func NewLicence(p price.Price, txnID string, creator admin.Creator) Licence {

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
		HintGrantMismatch:     false,
		StartDateUTC:          chrono.TimeUTCFrom(now),
		TrialStartUTC:         chrono.Date{},
		TrialEndUTC:           chrono.Date{},
		LatestTransactionID:   null.StringFrom(txnID),
		LatestPrice:           p,
		Status:                LicStatusAvailable,
		LatestInvitation:      InvitationJSON{},
		AssigneeID:            null.String{},
		RowTime:               admin.NewRowTime(),
	}
}

func (l Licence) Renewed(p price.Price, txnID string) Licence {
	now := time.Now()

	startTime := l.RenewalStartTime()
	period := dt.NewTimeRange(startTime).WithCycle(p.Cycle)

	l.CurrentPeriodStartUTC = chrono.TimeUTCFrom(period.Start)
	l.CurrentPeriodEndUTC = chrono.TimeUTCFrom(period.End)
	l.LatestPrice = p
	l.LatestTransactionID = null.StringFrom(txnID)
	l.UpdatedUTC = chrono.TimeUTCFrom(now)

	return l
}

func (l Licence) IsZero() bool {
	return l.ID == ""
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

// CreateInvitation builds a new invitation for this licence.
func (l Licence) CreateInvitation(params input.InvitationParams, claims admin.PassportClaims) (Licence, error) {
	inv, err := NewInvitation(params, claims)
	if err != nil {
		return Licence{}, err
	}

	l.HintGrantMismatch = false
	l.Status = LicStatusInvited
	l.LatestInvitation = InvitationJSON{inv}
	l.AssigneeID = null.String{} // Only exists after accepted.
	l.UpdatedUTC = chrono.TimeNow()

	return l, nil
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
	l.HintGrantMismatch = false
	l.Status = LicStatusAvailable
	l.LatestInvitation = InvitationJSON{}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeUTCNow()

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
