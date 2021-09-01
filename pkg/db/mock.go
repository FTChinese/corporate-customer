// +build !production

package db

import (
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/jmoiron/sqlx"
)

func MockMySQL() ReadWriteMyDBs {
	config.MustSetupViper(config.MustReadConfigFile())
	return MustNewMyDBs(false)
}

func MockTx() *sqlx.Tx {
	return MockMySQL().Write.MustBegin()
}
