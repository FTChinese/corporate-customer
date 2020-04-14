package sq

import "strings"

type Column struct {
	name     string
	alias    string
	ordering string
}

func NewColumn(name string) Column {
	return Column{
		name:  name,
		alias: "",
	}
}

func (c Column) AS(alias string) Column {
	c.alias = alias
	return c
}

func (c Column) Asc() Column {
	c.ordering = "ASC"

	return c
}

func (c Column) Desc() Column {
	c.ordering = "DESC"
	return c
}

func (c Column) Build() string {
	var buf strings.Builder

	buf.WriteString(c.name)

	if c.alias != "" {
		buf.WriteByte(' ')
		buf.WriteString("AS")
		buf.WriteByte(' ')
		buf.WriteString(c.alias)
		return buf.String()
	}

	if c.ordering != "" {
		buf.WriteByte(' ')
		buf.WriteString(c.ordering)
		return buf.String()
	}

	return buf.String()
}

func buildColumns(cols []Column) string {
	str := make([]string, len(cols))

	for i, c := range cols {
		str[i] = c.Build()
	}

	return strings.Join(str, ", ")
}
