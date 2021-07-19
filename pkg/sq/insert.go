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

// rowHolder builds `(?, ..., ?)` part of an insert statement, based on the number of columns.
func (i Insert) rowHolder() string {
	str := make([]string, len(i.cols))

	for j := range i.cols {
		str[j] = "?"
	}

	return "(" + strings.Join(str, ", ") + ")"
}

// placeholder builds the `(?, ..., ?), (?, ..., ?)` part of insert statement.
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

// InsertRow should be implemented by a type that can
// product an array of values that will be used
// as a row in SQL INSERT VALUES ().
// The values of all columns are put into an array.
type InsertRow interface {
	RowValues() []interface{}
}

// Enumerable should be implemented by an array of element
// to be used in a bulk insert.
type Enumerable interface {
	Each(handler func(row InsertRow))
}

// BuildInsertValues transform an array of InsertRow
// to the arg in sql's Exec method.
// Since each row forms an array, and each row's value
// is also an array, this operation flattens a 2-D array.
// It can then be used as a varidic.
func BuildInsertValues(rows Enumerable) []interface{} {
	var values = make([]interface{}, 0)

	rows.Each(func(row InsertRow) {
		values = append(values, row.RowValues()...)
	})

	return values
}
