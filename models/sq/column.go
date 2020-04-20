package sq

import "strings"

type colUsage int

const (
	colUsageSelect colUsage = iota
	colUsageUpdate
	colUsageOrdering
	colUsageComparison
)

// Column maps to SQL column name.
type Column struct {
	name     string
	alias    string
	ordering string
	value    string
	usage    colUsage
	operator string
}

func NewColumn(name string) Column {
	return Column{
		name:  name,
		alias: "",
	}
}

// AS set the alias in select statement.
// Example:
// ```
// NewColumn("user_name").AS("name")
// ```
// will produce: `user_name AS name`.
func (c Column) AS(alias string) Column {
	c.alias = alias
	c.usage = colUsageSelect
	return c
}

// SetTo uses column in a single row update statement:
// `UPDATE tbl SET col = val`
// NewColumn("col").SetTo("val")
func (c Column) SetTo(v string) Column {
	c.value = v
	c.usage = colUsageUpdate
	return c
}

// EqualTo set the value for this column.
// It is usually used in a SET clause, or WHERE clause.
// Example:
// For the following SQLX syntax
// ```
// UPDATE user.account
// SET user_name = :user_name
// WHERE id = :id
// ```
// NewColumn("id").EqualTo(":id")
func (c Column) EqualTo(v string) Column {
	c.value = v
	c.operator = "="
	c.usage = colUsageComparison
	return c
}

func (c Column) GreaterThan(v string) Column {
	c.value = v
	c.operator = ">"
	c.usage = colUsageComparison
	return c
}

func (c Column) GreaterOrEqual(v string) Column {
	c.value = v
	c.operator = ">="
	c.usage = colUsageComparison
	return c
}

func (c Column) LessThan(v string) Column {
	c.value = v
	c.operator = "<"
	c.usage = colUsageComparison
	return c
}

func (c Column) LessOrEqual(v string) Column {
	c.value = v
	c.operator = "<="
	c.usage = colUsageComparison
	return c
}

// Asc set ASCENDING order when used in ORDER By clause.
func (c Column) Asc() Column {
	c.ordering = "ASC"
	c.usage = colUsageOrdering
	return c
}

// Desc for DESCENDING order in ORDER BY clause.
func (c Column) Desc() Column {
	c.ordering = "DESC"
	c.usage = colUsageOrdering
	return c
}

// Build produces SQL string.
// AS is mutually exclusive with Asc or Desc.
func (c Column) Build() string {
	var buf strings.Builder

	buf.WriteString(c.name)

	switch c.usage {
	case colUsageSelect:
		if c.alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(c.alias)
		}
		return buf.String()

	case colUsageUpdate:
		buf.WriteString(" = ")
		buf.WriteString(c.value)
		return buf.String()

	case colUsageComparison:
		buf.WriteByte(' ')
		buf.WriteString(c.operator)
		buf.WriteByte(' ')
		buf.WriteString(c.value)
		return buf.String()

	case colUsageOrdering:
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
