package builder

import "strings"

type OrderBy struct {
	cols []Column
}

func NewOrderBy() OrderBy {
	return OrderBy{
		cols: make([]Column, 0),
	}
}

func (o OrderBy) AddColumn(col Column) OrderBy {
	o.cols = append(o.cols, col)

	return o
}

func (o OrderBy) Build() string {
	if len(o.cols) == 0 {
		return ""
	}

	var buf strings.Builder

	for _, v := range o.cols {
		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(' ')
		buf.WriteString(v.Build())
	}

	return "ORDER BY" + buf.String()
}
