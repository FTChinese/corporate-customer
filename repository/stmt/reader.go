package stmt

const readerAccountCols = `
u.user_id AS ftc_id,
u.email AS email,
u.user_name AS user_name,
IFNULL(u.is_vip, FALSE) AS is_vip`

const MembershipSelectCols = `
m.id AS subs_id,
m.vip_id AS subs_compound_id,
m.ftc_user_id AS subs_ftc_id,
m.wx_union_id AS subs_union_id,
m.vip_id_alias AS legacy_wx_id,
m.vip_type AS legacy_tier
m.expire_time AS legacy_expire
m.member_tier AS tier,
m.billing_cycle AS cycle,
m.expire_date AS expire_date,
m.auto_renewal AS auto_renew,
m.payment_method AS payment_method,
m.stripe_subscription_id AS stripe_subs_id,
m.stripe_plan_id AS stripe_plan_id,
m.sub_status AS subs_status,
m.apple_subscription_id AS apple_subs_id`

const SelectReader = `
SELECT + ` + readerAccountCols + `,
` + MembershipSelectCols + `
FROM cmstmp01.uerinfo AS u
	LEFT JOIN premium.ftc_vip AS m
	ON u.user_id = m.vip_id
WHERE u.email = ?
LIMIT 1`

// membershipUpsertCols list shared columns
// for both inserting and updating membership.
const membershipUpsertCols = `
vip_type = :legacy_tier,
expire_time = :legacy_expire,
member_tier = :tier,
billing_cycle = :cycle,
expire_date = :expire_date,
auto_renewal = :auto_renew,
payment_method = :payment_method,
stripe_subscription_id = :stripe_subs_id,
stripe_plan = :stripe_plan,
sub_status = :subs_status,
apple_subscription_id = :app_subs_id`

const InsertMembership = `
INSERT INTO premium.ftc_vip
SET id = :subs_id,
	vip_id = :subs_compound_id,
	ftc_user_id = :subs_ftc_id,
	wx_union_id = :subs_union_id,
	vip_id_alias = :legacy_wx_id,
` + membershipUpsertCols

// This statement updates an existing
// membership and conditionally set the
// id field is missing.
const UpdateMembership = `
UPDATE premium.ftc_vip
SET id = IFNULL(id, :subs_id),
` + membershipUpsertCols + `
WHERE vip_id = :subs_compound_id
LIMIT 1`
