package sq

import "testing"

func TestNewFrom(t *testing.T) {
	f1 := NewFrom(
		NewTable("b2b.licence").AS("l")).
		LeftJoin(
			NewTable("subs.plan").AS("p")).
		On("l.plan_id = p.id").
		LeftJoin(
			NewTable("cmstmp01.userinfo").AS("u")).
		On("l.assignee_id = u.user_id")

	t.Logf("From clause: %s", f1.Build())
}
