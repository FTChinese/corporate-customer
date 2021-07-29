package letter

import (
	"fmt"
	"github.com/FTChinese/go-rest/enum"
	"strconv"
	"strings"
	"text/template"
)

var tmplCache = map[string]*template.Template{}

var tierCN = map[enum.Tier]string{
	enum.TierStandard: "标准版",
	enum.TierPremium:  "高端版",
}

var funcMap = template.FuncMap{
	"formatFloat": func(f float64) string {
		return strconv.FormatFloat(f, 'f', 2, 32)
	},
	"currency": func(f float64) string {
		return fmt.Sprintf("¥ %.2f",
			f)
	},
	"tierSC": func(t enum.Tier) string {
		return tierCN[t]
	},
}

func Render(name string, ctx interface{}) (string, error) {
	tmpl, ok := tmplCache[name]
	var err error
	if !ok {
		tmplStr, ok := templates[name]
		if !ok {
			return "", fmt.Errorf("template %s not found", name)
		}

		tmpl, err = template.
			New(name).
			Funcs(funcMap).
			Parse(tmplStr)

		if err != nil {
			return "", err
		}
		tmplCache[name] = tmpl
	}

	var body strings.Builder
	err = tmpl.Execute(&body, ctx)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
