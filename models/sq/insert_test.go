package sq

import "testing"

func TestInsert_Build(t *testing.T) {
	insert := NewInsert().
		Into(NewTable("b2b.order")).
		AddColumn(NewColumn("id")).
		AddColumn(NewColumn("plan_id")).
		AddColumn(NewColumn("discount_id")).
		AddColumn(NewColumn("licence_id")).
		AddColumn(NewColumn("team_id")).
		AddColumn(NewColumn("checkout_id")).
		AddColumn(NewColumn("amount")).
		AddColumn(NewColumn("cycle_count")).
		AddColumn(NewColumn("trial_days")).
		AddColumn(NewColumn("kind")).
		AddColumn(NewColumn("created_utc")).
		Rows(5)

	t.Logf("Row holder: %s", insert.rowHolder())
	t.Logf("All row holder: %s", insert.placeholder())

	t.Logf("%s", insert.Build())
}
