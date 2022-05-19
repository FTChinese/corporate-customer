package pkg

type LegalTeaser struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
	Title  string `json:"title"`
}

type LegalDoc struct {
	LegalTeaser
	Body string `json:"body"`
}

type LegalList struct {
	PagedList
	Data []LegalTeaser `json:"data"`
}
