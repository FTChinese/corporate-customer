package licence

const StmtVersioned = `
INSERT INTO b2b.licence_version
SET action_kind = :action_kind,
	ante_change = :ante_change,
	membership_version_id = :membership_version_id,
	mismatched_member = :mismatched_member,
	post_change = :post_change,
	created_utc = :created_utc
`
