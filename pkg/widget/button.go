package widget

import (
	"html/template"
	"strings"
)

const btnTmplText = `
<button
	class="btn{{ if .Block }} btn-block{{ end }} btn-primary"
	type="submit"
	data-disable-with="{{ .DisableWith }}">{{ .Text }}</button>`

var btnTmpl = template.Must(template.New("button").Parse(btnTmplText))

type BtnType string

const (
	BtnTypeButton = "button"
	BtnTypeSubmit = "submit"
)

type BtnStyle string

const (
	BtnStylePrimary   = "primary"
	BtnStyleSecondary = "secondary"
	BtnStyleSuccess   = "success"
	BtnStyleDanger    = "danger"
	BtnStyleWarning   = "warning"
	BtnStyleInfo      = "info"
	BtnStyleLight     = "light"
	BtnStyleDark      = "dark"
	BtnStyleLink      = "link"
)

func (s BtnStyle) Value(outline bool) string {
	var base = "btn-"
	if outline {
		base = base + "outline-"
	}
	return base + string(s)
}

type BtnSize string

const (
	BtnSizeSmall = "sm"
	BtnSizeLarge = "lg"
	BtnSizeBlock = "block"
)

func (s BtnSize) Value() string {
	return "btn-" + string(s)
}

type Button struct {
	Text        string
	DisableWith string
	ClassName   []string
	Type        BtnType
	Style       BtnStyle
	Size        BtnSize
	Outline     bool
}

var PrimaryBtn = Button{
	Text:        "",
	DisableWith: "",
	ClassName:   []string{"btn"},
	Type:        BtnTypeSubmit,
	Style:       BtnStylePrimary,
	Size:        "",
	Outline:     false,
}

var PrimaryBlockBtn = Button{
	Text:        "",
	DisableWith: "",
	ClassName:   []string{"btn"},
	Type:        BtnTypeSubmit,
	Style:       BtnStylePrimary,
	Size:        BtnSizeBlock,
	Outline:     false,
}

func (b Button) SetName(n string) Button {
	b.Text = n
	return b
}

func (b Button) SetDisabledText(t string) Button {
	b.DisableWith = t
	return b
}

func (b Button) SetType(t BtnType) Button {
	b.Type = t
	return b
}

func (b Button) SetStyle(s BtnStyle) Button {
	b.Style = s
	return b
}

func (b Button) SetSize(s BtnSize) Button {
	b.Size = s
	return b
}

func (b Button) SetOutline() Button {
	b.Outline = true
	return b
}

func (b Button) AddClass(c string) Button {
	b.ClassName = append(b.ClassName, c)
	return b
}

// ClassName builds the html tag class attribute.
func (b Button) GetClassName() string {

	if b.Style != "" {
		b.ClassName = append(b.ClassName, b.Style.Value(b.Outline))
	}

	if b.Size != "" {
		b.ClassName = append(b.ClassName, b.Size.Value())
	}

	return strings.Join(b.ClassName, " ")
}
