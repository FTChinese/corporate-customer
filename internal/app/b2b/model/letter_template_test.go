package model

import "testing"

func TestCompileTemplate(t *testing.T) {
	t.Logf("%s", tmpl.DefinedTemplates())
	t.Logf("Template name: %s", tmpl.Name())
}
