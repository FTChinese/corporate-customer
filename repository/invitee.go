package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/reader"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

// LoadInvitation
func (env Env) FindInvitation(token string) (admin.ExpandedInvitation, error) {
	var i admin.ExpandedInvitation

	err := env.db.Get(&i, stmt.FindExpandedInvitation, token)

	if err != nil {
		return admin.ExpandedInvitation{}, err
	}

	return i, nil
}

// FindLicence tries to retrieve an expanded
// licence by id after an invitation is verified.
func (env Env) FindLicence(id string) (admin.ExpandedLicence, error) {
	var ls admin.LicenceSchema

	err := env.db.Get(&ls, stmt.FindExpandedLicence, id)

	if err != nil {
		return admin.ExpandedLicence{}, err
	}

	return ls.ExpandedLicence(), err
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

func (env Env) LoadInvitee(email string) (reader.Invitee, error) {
	var i reader.Invitee
	if err := env.db.Get(&i, stmtInvitee, email); err != nil {
		return reader.Invitee{}, err
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

func (env Env) CreateReader(s reader.SignUp) error {
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
func (env Env) TakeSnapshot(snp reader.MemberSnapshot) error {
	_, err := env.db.NamedExec(stmtTakeSnapshot, snp)

	if err != nil {
		return err
	}

	return nil
}

// GrantLicence grants a licence to a reader.
func (env Env) GrantLicence(expInv admin.ExpandedInvitation, expLic admin.ExpandedLicence) error {
	tx, err := env.beginGrantTx()
	if err != nil {
		return err
	}

	_, err = tx.LockInvitation(expInv.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.LockLicence(expLic.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	mmb, err := tx.LockMembership(expInv.Assignee.FtcID.String)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	newMmb := mmb.BuildOn(expLic)

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
			_ = env.TakeSnapshot(reader.MemberSnapshot{
				SnapshotID: reader.GenerateSnapshotID(),
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
