package pkg

type SQLWhere struct {
	Clause string
	Values []interface{}
}

func (w SQLWhere) AddValues(v ...interface{}) SQLWhere {
	w.Values = append(w.Values, v...)
	return w
}
