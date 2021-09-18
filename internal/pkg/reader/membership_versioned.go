package reader

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type ArchiveName string

const (
	ArchiveNameB2B ArchiveName = "b2b"
)

type ArchiveAction string

const (
	ArchiveActionGrant  ArchiveAction = "grant"
	ArchiveActionRenew  ArchiveAction = "renew"
	ArchiveActionRevoke ArchiveAction = "revoke"
)

type Archiver struct {
	Name   ArchiveName
	Action ArchiveAction
}

func (a Archiver) String() string {
	return fmt.Sprintf("%s.%s", a.Name, a.Action)
}

func B2BArchiver(a ArchiveAction) Archiver {
	return Archiver{
		Name:   ArchiveNameB2B,
		Action: a,
	}
}

type MembershipJSON struct {
	Membership
}

// Value implements Valuer interface by saving the entire
// type as JSON string, or null if it is a zero value.
func (m MembershipJSON) Value() (driver.Value, error) {
	// For zero value, save as NULL.
	if m.IsZero() {
		return nil, nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (m *MembershipJSON) Scan(src interface{}) error {
	// Handle null value.
	if src == nil {
		*m = MembershipJSON{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp MembershipJSON
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*m = tmp
		return nil

	default:
		return errors.New("incompatible type to scna to MembershipJSON")
	}
}

const StmtVersionMembership = `
INSERT INTO premium.member_version
SET id = :snapshot_id,
	ante_change = :ante_change,
	created_by = :created_by,
	created_utc = :created_utc,
	b2b_transaction_id = :b2b_transaction_id,
	post_change = :post_change,
	retail_order_id = :retail_order_id,
`

// MembershipVersioned stores a specific version of membership.
// Since membership is constantly changing, we keep all
// versions of modification in a dedicated table.
type MembershipVersioned struct {
	ID               string         `json:"id" db:"snapshot_id"`
	AnteChange       MembershipJSON `json:"anteChange" db:"ante_change"` // Membership before being changed
	CreatedBy        null.String    `json:"createdBy" db:"created_by"`
	CreatedUTC       chrono.Time    `json:"createdUtc" db:"created_utc"`
	B2BTransactionID null.String    `json:"b2bTransactionId" db:"b2b_transaction_id"`
	PostChange       MembershipJSON `json:"postChange" db:"post_change"`       // Membership after being changed.
	RetailOrderID    null.String    `json:"retailOderId" db:"retail_order_id"` // Only exists when user is performing renewal or upgrading.
}

// IsZero tests if a snapshot exists.
func (s MembershipVersioned) IsZero() bool {
	return s.ID == ""
}

func (s MembershipVersioned) WithB2BTxnID(id string) MembershipVersioned {
	s.B2BTransactionID = null.StringFrom(id)
	return s
}

// WithRetailOrderID sets the retail order id when taking a
// snapshot.
func (s MembershipVersioned) WithRetailOrderID(id string) MembershipVersioned {
	s.RetailOrderID = null.StringFrom(id)
	return s
}

// WithPriorVersion stores a previous version of membership
// before modification.
// It may not exist for newly created membership.
func (s MembershipVersioned) WithPriorVersion(m Membership) MembershipVersioned {
	if m.IsZero() {
		return s
	}

	s.AnteChange = MembershipJSON{m}

	return s
}

// Version takes a snapshots of membership both before and after modification.
func (m Membership) Version(by Archiver) MembershipVersioned {
	if m.IsZero() {
		return MembershipVersioned{}
	}

	return MembershipVersioned{
		ID:               pkg.SnapshotID(),
		AnteChange:       MembershipJSON{}, // Optional. Only exists if exists prior to the latest one.
		CreatedBy:        null.StringFrom(by.String()),
		CreatedUTC:       chrono.TimeNow(),
		B2BTransactionID: null.String{},
		PostChange:       MembershipJSON{m},
		RetailOrderID:    null.String{},
	}
}
