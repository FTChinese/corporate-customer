package sq

import (
	"testing"
)

func TestNewOrderBy(t *testing.T) {
	ob := NewOrderBy().
		AddColumn(NewColumn("p.tier").Asc()).
		AddColumn(NewColumn("p.cycle").Desc())

	t.Logf("Order by: %s", ob.Build())
}
