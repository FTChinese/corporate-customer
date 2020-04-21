package widget

type FlashKind string

const (
	FlashKindSuccess = "alert-success"
	FlashKindDanger  = "alert-danger"
)

type Flash struct {
	Message     string
	Kind        FlashKind
	Dismissible bool
}

func NewFlash(msg string) *Flash {
	return &Flash{
		Message: msg,
	}
}

func (f *Flash) SetSuccess() *Flash {
	f.Kind = FlashKindSuccess
	return f
}

func (f *Flash) SetDanger() *Flash {
	f.Kind = FlashKindDanger
	return f
}

func (f *Flash) SetDismissible() *Flash {
	f.Dismissible = true
	return f
}
