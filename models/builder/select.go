package builder

import (
	"strconv"
	"strings"
)

type Select struct {
	rawCols []string
	from    string
	where   string
	orderBy string
	limit   int
	paged   bool
	lock    bool
}

func NewSelect() Select {
	return Select{
		rawCols: make([]string, 0),
		from:    "",
		orderBy: "",
		limit:   0,
		paged:   false,
	}
}

func (s Select) AddRawColumn(col string) Select {
	s.rawCols = append(s.rawCols, col)
	return s
}

func (s Select) From(from string) Select {
	s.from = from
	return s
}

func (s Select) Where(w string) Select {
	s.where = w
	return s
}

func (s Select) OrderBy(o string) Select {
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
	buf.WriteString(strings.Join(s.rawCols, ","))
	buf.WriteByte(' ')
	buf.WriteString("FROM ")
	buf.WriteString(s.from)
	if s.where != "" {
		buf.WriteByte(' ')
		buf.WriteString("WHERE ")
		buf.WriteString(s.where)
	}

	if s.orderBy != "" {
		buf.WriteByte(' ')
		buf.WriteString("ORDER BY ")
		buf.WriteString(s.orderBy)
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
