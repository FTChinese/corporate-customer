//go:build !production
// +build !production

package db

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/jmoiron/sqlx"
)

func MockMySQL() ReadWriteMyDBs {
	faker.MustSetupViper()
	return MustNewMyDBs()
}

func MockTx() *sqlx.Tx {
	return MockMySQL().Write.MustBegin()
}
