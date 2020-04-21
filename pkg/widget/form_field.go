package widget

type ControlType int

const (
	ControlTypeHidden ControlType = iota
	ControlTypeText
	ControlTypeEmail
	ControlTypePassword
	ControlTypeURL
	ControlTypeNumber
	ControlTypeSearch
	ControlTypeMonth
	ControlTypeWeek
	ControlTypeDate
	ControlTypeTel
	ControlTypeImage
	ControlTypeFile
	ControlTypeTextArea
)

var ControlTypeNames = [...]string{
	"hidden",
	"text",
	"email",
	"password",
	"url",
	"number",
	"search",
	"month",
	"date",
	"week",
	"tel",
	"image",
	"file",
	"textarea",
}

func (i ControlType) String() string {
	if i < 0 || i > ControlTypeTextArea {
		return ""
	}

	return ControlTypeNames[i]
}

func (i ControlType) IsTextArea() bool {
	return i == ControlTypeTextArea
}

type FormControl struct {
	Label       string
	ID          string
	Type        ControlType
	Name        string
	Value       string
	Placeholder string
	Required    bool
	ReadOnly    bool
	Checked     bool
	MinLength   int
	MaxLength   int
	Pattern     string
	Desc        string
	ErrMsg      string
}

func NewFormControl(t ControlType) FormControl {
	return FormControl{
		Type: t,
	}
}
