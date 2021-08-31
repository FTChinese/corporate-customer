package sq

type Where struct {
}

func NewWhere(col Column) Where {
	return Where{}
}

func NewWhereParen(p Paren) Where {
	return Where{}
}

func (w Where) AddColumn(col Column) Where {
	return w
}

func (w Where) AddParen(p Paren) Where {
	return w
}

func (w Where) And(col Column) Where {
	return w
}

func (w Where) AndParen(p Paren) Where {
	return w
}

func (w Where) Or(col Column) Where {
	return w
}

func (w Where) OrParen(p Paren) Where {
	return w
}
