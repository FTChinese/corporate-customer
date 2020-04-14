package sq

import (
	"testing"
)

func TestColumn_AS(t *testing.T) {
	col := NewColumn("l.id").AS("licence_id")

	t.Logf("AS clause: %s", col.Build())
}

func TestColumn_Asc(t *testing.T) {
	col := NewColumn("l.created_utc").Asc()

	t.Logf("ORDER BY ... ASC clause: %s", col.Build())
}

func TestColumn_Desc(t *testing.T) {
	col := NewColumn("l.created_utc").Desc()

	t.Logf("ORDER BY ... DESC clause: %s", col.Build())
}
