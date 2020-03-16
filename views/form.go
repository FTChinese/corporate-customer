package views

type Link struct {
	Text string
	Href string
}

type Button struct {
	DisableWith string
	Text        string
}

type Form struct {
	Disabled  bool
	Action    string
	Fields    []FormField
	SubmitBtn Button
	CancelBtn Link
	DeleteBtn Link
}

func (f Form) WithErrors(e map[string]string) Form {
	for i, v := range f.Fields {
		f.Fields[i].ErrMsg = e[v.Name]
	}

	return f
}
