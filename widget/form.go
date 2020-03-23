package widget

type Form struct {
	Disabled  bool
	Action    string
	Fields    []FormControl
	SubmitBtn Button
	CancelBtn Link
	DeleteBtn Link
}

func NewForm(action string) Form {
	return Form{
		Action: action,
		Fields: []FormControl{},
	}
}

func (f Form) Disable() Form {
	f.Disabled = true
	return f
}

func (f Form) AddControl(c FormControl) Form {
	f.Fields = append(f.Fields, c)
	return f
}

func (f Form) WithErrors(e map[string]string) Form {
	for i, v := range f.Fields {
		f.Fields[i].ErrMsg = e[v.Name]
	}

	return f
}
