package sq

type Paren struct {
	cols []Column
}

func NewParen() Paren {
	return Paren{
		cols: []Column{},
	}
}

func (p Paren) And(a Column, b Column) Paren {

	return p
}

func (p Paren) Or(a Column, b Column) Paren {
	return p
}
