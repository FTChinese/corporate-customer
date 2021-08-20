package sq

import "strings"

// BulkInsert build INSERT statement like:
// INSERT INTO tbl (
//    col1,
//    col2,
// ) VALUES
// (?, ?),
// (?, ?),
// (?, ?)
type BulkInsert struct {
	tbl      Table
	cols     []Column
	rowCount int
}

func NewBulkInsert() BulkInsert {
	return BulkInsert{}
}

func (i BulkInsert) Into(t Table) BulkInsert {
	i.tbl = t
	return i
}

func (i BulkInsert) SetColumns(cols ...Column) BulkInsert {
	i.cols = cols
	return i
}

func (i BulkInsert) AddColumn(col Column) BulkInsert {
	i.cols = append(i.cols, col)
	return i
}

func (i BulkInsert) Rows(n int) BulkInsert {
	i.rowCount = n
	return i
}

// rowHolder builds `(?, ..., ?)` part of an insert statement, based on the number of columns.
func (i BulkInsert) rowHolder() string {
	str := make([]string, len(i.cols))

	for j := range i.cols {
		str[j] = "?"
	}

	return "(" + strings.Join(str, ", ") + ")"
}

// placeholder builds the `(?, ..., ?), (?, ..., ?)` part of insert statement.
func (i BulkInsert) placeholder() string {
	holder := i.rowHolder()

	str := make([]string, i.rowCount)

	for j := 0; j < i.rowCount; j++ {
		str[j] = holder
	}

	return strings.Join(str, ", ")
}

func (i BulkInsert) Build() string {
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

// BulkInsertRow should be implemented by a type that can
// product an array of values that will be used
// as a row in SQL INSERT VALUES ().
// The values of all columns are put into an array.
type BulkInsertRow interface {
	RowValues() []interface{}
}

// Enumerable should be implemented by an array of element
// to be used in a bulk insert.
// Each function is passed each element of the array.
type Enumerable interface {
	Each(handler func(row BulkInsertRow))
}

// BuildBulkInsertValues transform an array of BulkInsertRow
// to the arg in sql's Exec method.
// Since each row forms an array, and each row's value
// is also an array, this operation flattens a 2-D array.
// It can then be used as a varidic.
// Example:
// ```
// type OrderItems []OrderItem
// func (ci OrderItems) Each(handler func(row BulkInsertRow)) {
//	for _, c := range ci {
//		handler(c)
//	}
// }
// ```
// Why not use the array to build values directly?
// For example, if you have one like:
// func ConcatBulkInsertValues(rows []BulkInsertRow) []interface{}
// Pass a slice of types implementing the interface if not
// allowed since the the slice of implementer is not a slice
// of BulkInsertRow.
func BuildBulkInsertValues(rows Enumerable) []interface{} {
	var values = make([]interface{}, 0)

	// Concatenate each row's values.
	rows.Each(func(row BulkInsertRow) {
		values = append(values, row.RowValues()...)
	})

	return values
}
