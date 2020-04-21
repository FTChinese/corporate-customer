package sq

import "testing"

func TestSelect_Build(t *testing.T) {
	s := NewSelect().
		AddColumn(NewColumn("l.id").AS("licence_id")).
		AddColumn(NewColumn("l.team_id").AS("team_id")).
		AddColumn(NewColumn("l.expire_date").AS("expire_date")).
		AddColumn(NewColumn("l.assignee_id").AS("assignee_id")).
		AddColumn(NewColumn("l.is_active").AS("is_active")).
		AddColumn(NewColumn("l.created_utc").AS("created_utc")).
		AddColumn(NewColumn("l.updated_utc").AS("updated_utc")).
		AddColumn(NewColumn("p.id").AS("plan_id")).
		AddColumn(NewColumn("p.price").AS("price")).
		AddColumn(NewColumn("p.tier").AS("tier")).
		AddColumn(NewColumn("p.cycle").AS("cycle")).
		AddColumn(NewColumn("p.trial_days").AS("trial_days")).
		AddColumn(NewColumn("u.user_id").AS("ftc_id")).
		AddColumn(NewColumn("u.email").AS("email")).
		AddColumn(NewColumn("u.user_name").AS("user_name")).
		AddColumn(NewColumn("u.is_vip").AS("is_vip")).
		From(
			NewFrom(NewTable("b2b.licence").AS("l")).
				LeftJoin(NewTable("subs.plan").AS("p")).
				On("l.plan_id = p.id").
				LeftJoin(NewTable("cmstmp01.userinfo").AS("u")).
				On("l.assignee_id = u.user_id"),
		).
		Where("l.id = ? AND l.team_id = ?").
		Limit(1).
		Build()

	t.Logf("SQL string: %s", s)
}
