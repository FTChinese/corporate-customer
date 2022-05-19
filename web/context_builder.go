package web

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/flosch/pongo2/v4"
	"time"
)

type ContextBuilder struct {
	version string
	error   string
	title   string
	ctx     pongo2.Context
}

func NewContextBuilder(version string) *ContextBuilder {
	return &ContextBuilder{
		version: version,
		ctx:     pongo2.Context{},
	}
}

func (cb *ContextBuilder) SetErr(msg string) *ContextBuilder {
	cb.error = msg
	return cb
}

func (cb *ContextBuilder) SetTitle(t string) *ContextBuilder {
	cb.title = t
	return cb
}

func (cb *ContextBuilder) SetLegalList(list []pkg.LegalTeaser) *ContextBuilder {
	cb.ctx["legalTeasers"] = list
	return cb
}

func (cb *ContextBuilder) SetLegalDoc(doc pkg.LegalDoc) *ContextBuilder {
	cb.ctx["legalDoc"] = doc
	return cb
}

func (cb *ContextBuilder) Build() pongo2.Context {
	cb.ctx["head"] = DefaultHead
	cb.ctx["title"] = cb.title
	cb.ctx["footer"] = Footer{
		Matrix:        DefaultFooterCols,
		Year:          time.Now().Year(),
		ClientVersion: "",
		ServerVersion: cb.version,
	}
	cb.ctx["error"] = cb.error

	return cb.ctx
}
