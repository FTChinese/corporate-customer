package builder

import "strings"

type Table struct {
	name  string
	alias string
}

func NewTable(name string) Table {
	return Table{
		name: name,
	}
}

func (t Table) AS(alias string) Table {
	t.alias = alias
	return t
}

func (t Table) Build() string {
	var buf strings.Builder
	buf.WriteString(t.name)

	if t.alias != "" {
		buf.WriteString(" AS ")
		buf.WriteString(t.alias)
	}

	return buf.String()
}
