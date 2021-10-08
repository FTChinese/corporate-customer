package web

import "testing"

func TestEmbedFS(t *testing.T) {
	entries, err := templates.ReadDir(".")
	if err != nil {
		t.Error(err)
	}

	for _, entry := range entries {
		t.Logf("Name %s", entry.Name())
	}

	b, err := templates.ReadFile("template/b2b/home.html")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", b)
}
