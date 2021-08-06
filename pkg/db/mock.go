// +build !production

package db

import (
	"github.com/FTChinese/ftacademy/pkg/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

func readConfigFile() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filepath.Join(home, "config", "api.toml"))
}

func MockMySQL() ReadWriteMyDBs {
	b, err := readConfigFile()
	if err != nil {
		panic(err)
	}

	config.MustSetupViper(b)
	return MustNewMyDBs(false)
}
