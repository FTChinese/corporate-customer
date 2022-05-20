package pkg

import "github.com/gomarkdown/markdown"

type LegalTeaser struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
	Title  string `json:"title"`
}

type LegalDoc struct {
	LegalTeaser
	Body string `json:"body"`
}

func (l LegalDoc) Rendered() LegalDoc {
	output := markdown.ToHTML([]byte(l.Body), nil, nil)
	return LegalDoc{
		LegalTeaser: l.LegalTeaser,
		Body:        string(output),
	}
}

type LegalList struct {
	PagedList
	Data []LegalTeaser `json:"data"`
}
