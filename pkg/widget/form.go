package widget

type Method string

const (
	MethodPost = "post"
	MethodGet  = "get"
)

type Form struct {
	Disabled  bool
	Method    Method
	Action    string
	Fields    []FormControl
	SubmitBtn Button
	CancelBtn Link
	DeleteBtn Link
}

func NewForm(action string) Form {
	return Form{
		Action: action,
		Method: MethodPost,
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

func (f Form) SetMethod(m Method) Form {
	f.Method = m

	return f
}

func (f Form) WithErrors(e map[string]string) Form {
	for i, v := range f.Fields {
		f.Fields[i].ErrMsg = e[v.Name]
	}

	return f
}
