package builder

import "strings"

type From struct {
	parts []string
}

func NewFrom(t Table) From {
	return From{
		parts: []string{
			t.Build(),
		},
	}
}

func (f From) Join(t Table) From {
	f.parts = append(f.parts, "JOIN "+t.Build())

	return f
}

func (f From) LeftJoin(t Table) From {
	f.parts = append(f.parts, "LEFT JOIN "+t.Build())

	return f
}

func (f From) On(condition string) From {
	f.parts = append(f.parts, "ON "+condition)

	return f
}

func (f From) Build() string {
	return strings.Join(f.parts, " ")
}
