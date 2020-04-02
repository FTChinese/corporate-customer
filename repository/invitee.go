package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
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
	u.is_vip AS is_vip
	m.id AS mmb_id,
	m.member_tier AS mmb_tier,
	m.billing_cycle AS mmb_cycle,
	m.expire_date AS mm_expire_date,
	m.payment_method AS mmb_pay_method,
FROM cmstmp01.uerinfo AS u
	LEFT JOIN premium.ftc_vip AS m
	ON u.user_id = m.vip_id
WHERE u.email = ?
LIMIT 1`

func (env Env) LoadInvitee(email string) (admin.Invitation, error) {
	var i admin.Invitee
	if err := env.db.Get(&i, stmtInvitee, email); err != nil {
		return admin.Invitation{}, err
	}

	return i, nil
}
