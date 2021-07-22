// +build !production

package db

func MockMySQL() ReadWriteMyDBs {
	return MustNewMyDBs(false)
}
