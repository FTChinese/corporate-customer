package sq

import "strings"

// Insert build INSERT statement like:
// INSERT INTO tbl (
//    col1,
//    col2,
// ) VALUES
// (?, ?),
// (?, ?),
// (?, ?)
type Insert struct {
	tbl      Table
	cols     []Column
	rowCount int
}

func NewInsert() Insert {
	return Insert{}
}

func (i Insert) Into(t Table) Insert {
	i.tbl = t
	return i
}

func (i Insert) SetColumns(cols []Column) Insert {
	i.cols = cols
	return i
}

func (i Insert) AddColumn(col Column) Insert {
	i.cols = append(i.cols, col)
	return i
}

func (i Insert) Rows(n int) Insert {
	i.rowCount = n
	return i
}

func (i Insert) rowHolder() string {
	str := make([]string, len(i.cols))

	for j := range i.cols {
		str[j] = "?"
	}

	return "(" + strings.Join(str, ", ") + ")"
}

func (i Insert) placeholder() string {
	holder := i.rowHolder()

	str := make([]string, i.rowCount)

	for j := 0; j < i.rowCount; j++ {
		str[j] = holder
	}

	return strings.Join(str, ", ")
}

func (i Insert) Build() string {
	var buf strings.Builder
	buf.WriteString("INSERT INTO ")
	buf.WriteString(i.tbl.Build())
	buf.WriteByte(' ')
	buf.WriteByte('(')
	buf.WriteString(buildColumns(i.cols))
	buf.WriteByte(')')
	buf.WriteString(" VALUES ")
	buf.WriteString(i.placeholder())

	return buf.String()
}
