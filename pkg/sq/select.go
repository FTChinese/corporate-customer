package sq

import (
	"strconv"
	"strings"
)

type Select struct {
	columns []Column
	from    From
	where   string
	orderBy OrderBy
	limit   int
	paged   bool
	lock    bool
}

func NewSelect() Select {
	return Select{
		columns: make([]Column, 0),
	}
}

func (s Select) SetColumns(cols []Column) Select {
	s.columns = cols
	return s
}

func (s Select) AddColumn(col Column) Select {
	s.columns = append(s.columns, col)
	return s
}

func (s Select) From(from From) Select {
	s.from = from
	return s
}

func (s Select) Where(w string) Select {
	s.where = w
	return s
}

func (s Select) OrderBy(o OrderBy) Select {
	s.orderBy = o
	return s
}

func (s Select) Limit(rows int) Select {
	s.limit = rows
	return s
}

func (s Select) Paged() Select {
	s.paged = true
	return s
}

func (s Select) Lock() Select {
	s.lock = true
	return s
}

// Build produces a SQL SELECT statement:
// SELECT ...
// FROM ...
// WHERE ...
// ORDER BY ...
// LIMIT ? [OFFSET ?]
func (s Select) Build() string {
	var buf strings.Builder

	buf.WriteString("SELECT ")
	buf.WriteString(buildColumns(s.columns))

	buf.WriteByte(' ')

	buf.WriteString(s.from.Build())

	if s.where != "" {
		buf.WriteByte(' ')
		buf.WriteString("WHERE ")
		buf.WriteString(s.where)
	}

	if len(s.orderBy.cols) != 0 {
		buf.WriteByte(' ')
		buf.WriteString(s.orderBy.Build())
	}

	if s.limit > 0 {
		buf.WriteByte(' ')
		buf.WriteString("LIMIT ")
		buf.WriteString(strconv.Itoa(s.limit))

	} else if s.paged {
		buf.WriteByte(' ')
		buf.WriteString("LIMIT ? OFFSET ?")
	}

	if s.lock {
		buf.WriteByte(' ')
		buf.WriteString("FOR UPDATE")
	}

	return buf.String()
}
