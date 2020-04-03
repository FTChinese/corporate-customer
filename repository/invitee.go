package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/invitee"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

const stmtVerifyInvitation = stmt.Invitation + `
WHERE id = ?
LIMIT 1`

func (env Env) VerifyInvitation(id string) (admin.Invitation, error) {
	var i admin.Invitation

	err := env.db.Get(&i, stmtVerifyInvitation, id)

	if err != nil {
		return i, err
	}

	return i, nil
}

const stmtVerifyLicence = stmt.Licence + `
WHERE id = ?
LIMIT 1`

func (env Env) VerifyLicence(id string) (admin.Licence, error) {
	var l admin.Licence

	err := env.db.Get(&l, stmtVerifyLicence, id)

	if err != nil {
		return l, err
	}

	return l, err
}

const stmtInvitee = `
SELECT u.user_id AS ftc_id,
	u.email AS email,
	u.is_vip AS is_vip,` +
	stmt.MembershipSelectCols + `
FROM cmstmp01.uerinfo AS u
	LEFT JOIN premium.ftc_vip AS m
	ON u.user_id = m.vip_id
WHERE u.email = ?
LIMIT 1`

func (env Env) LoadInvitee(email string) (invitee.Invitee, error) {
	var i invitee.Invitee
	if err := env.db.Get(&i, stmtInvitee, email); err != nil {
		return invitee.Invitee{}, err
	}

	return i, nil
}

const stmtCreateReader = `
INSERT INTO cmstmp01.userinfo
SET user_id = :ftc_id,
	email = :email,
	password = MD5(:password),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const stmtCreateReaderProfile = `
INSERT INTO user_db.profile
SET user_id = :ftc_id`

const stmtSaveReaderVrf = `
INSERT INTO user_db.email_verify
SET ftc_id = :ftc_id,
	email = :email,
	token = UNHEX(:token),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

func (env Env) CreateReader(s invitee.SignUp) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(stmtCreateReader, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(stmtCreateReaderProfile, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(stmtSaveReaderVrf, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

const stmtTakeSnapshot = `
INSERT INTO premium.member_snapshot
SET id = :snapshot_id,
	reason = :reason,
	created_utc = UTC_TIMESTAMP(),
	member_id = :subs_id,
	compound_id = :subs_compound_id,
	ftc_user_id = :subs_ftc_id,
	wx_union_id = :subs_union_id,
	tier = :tier,
	cycle = :cycle,
	expire_date = :expire_date,
	auto_renewal = :auto_renew,
	payment_method = :payment_method,
	stripe_subscription_id = :stripe_subs_id,
	stripe_plan = :stripe_plan,
	sub_status = :subs_status,
	apple_subscription_id = :app_subs_id`

// TakeSnapshot backs up a membership before
// modifying it.
func (env Env) TakeSnapshot(snp invitee.MemberSnapshot) error {
	_, err := env.db.NamedExec(stmtTakeSnapshot, snp)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) GrantLicence(expInv admin.ExpandedInvitation) error {
	tx, err := env.beginGrantTx()
	if err != nil {
		return err
	}

	inv, err := tx.LockInvitation(expInv.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.LockLicence(inv.LicenceID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	mmb, err := tx.LockMembership(expInv.Licence.AssigneeID.String)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	newMmb := mmb.BuildOn(expInv.Licence)

	// Create new membership based on licence
	if mmb.IsZero() {
		err := tx.InsertMembership(newMmb)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		// Update current membership based on
		// licence.
		err := tx.UpdateMembership(newMmb)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		// Back up.
		go func() {
			_ = env.TakeSnapshot(invitee.MemberSnapshot{
				SnapshotID: invitee.GenerateSnapshotID(),
				Reason:     enum.SnapshotReasonB2B,
				CreatedUTC: chrono.TimeNow(),
				Membership: mmb,
			})
		}()
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
