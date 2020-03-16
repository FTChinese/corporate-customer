package views

type InputType int

const (
	InputTypeHidden InputType = iota
	InputTypeText
	InputTypeEmail
	InputTypePassword
	InputTypeURL
	InputTypeNumber
	InputTypeSearch
	InputTypeMonth
	InputTypeWeek
	InputTypeDate
	InputTypeTel
	InputTypeImage
	InputTypeFile
	InputTypeTextArea
)

var inputTypeNames = [...]string{
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

func (i InputType) String() string {
	if i < 0 || i > InputTypeTextArea {
		return ""
	}

	return inputTypeNames[i]
}

func (i InputType) IsTextArea() bool {
	return i == InputTypeTextArea
}

type FormField struct {
	Label       string
	ID          string
	Type        InputType
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
