// +build !production

package db

import (
	"github.com/FTChinese/ftacademy/pkg/config"
)

func MockMySQL() ReadWriteMyDBs {
	config.MustSetupViper(config.MustReadConfigFile())
	return MustNewMyDBs(false)
}
