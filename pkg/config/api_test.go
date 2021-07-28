package config

import (
	"testing"
)

func TestLoadAPIConfig(t *testing.T) {
	MustSetupViper()
	a, err := LoadAPIConfig("api_keys.ftacademy")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", a)
}
