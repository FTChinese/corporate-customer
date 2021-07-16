package reader

import (
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

const StmtSaveSnapshot = `
INSERT INTO premium.member_snapshot
SET id = :snapshot_id,
	created_by = :created_by,
	created_utc = UTC_TIMESTAMP(),
	order_id = :order_id,
	compound_id = :compound_id,
	ftc_user_id = :ftc_id,
	wx_union_id = :union_id,
	tier = :tier,
	cycle = :cycle,
` + mUpsertSharedCols

type MemberSnapshot struct {
	SnapshotID string      `json:"id" db:"snapshot_id"`
	CreatedBy  null.String `json:"createdBy" db:"created_by"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	OrderID    null.String `json:"orderId" db:"order_id"` // Only exists when user is performing renewal or upgrading.
	Membership
}

// Snapshot takes a snapshot of membership, usually before modifying it.
func (m Membership) Snapshot(by Archiver) MemberSnapshot {
	if m.IsZero() {
		return MemberSnapshot{}
	}

	return MemberSnapshot{
		SnapshotID: pkg.SnapshotID(),
		CreatedBy:  null.StringFrom(by.String()),
		CreatedUTC: chrono.TimeNow(),
		Membership: m,
	}
}
