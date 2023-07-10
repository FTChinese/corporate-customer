package db

import (
	"fmt"
	"log"
	"time"

	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewMySQL creates opens a new sql connection.
func NewMySQL(c config.Connect) (*sqlx.DB, error) {
	cfg := &mysql.Config{
		User:   c.User,
		Passwd: c.Pass,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Params: map[string]string{
			"time_zone": `'+00:00'`,
		},
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// When connecting to production server it throws error:
	// packets.go:36: unexpected EOF
	//
	// See https://github.com/go-sql-driver/mysql/issues/674
	db.SetConnMaxLifetime(time.Second)
	log.Printf("Connected to MySQL %s", c.Host)
	return db, nil
}

func MustNewMySQL(c config.Connect) *sqlx.DB {
	db, err := NewMySQL(c)
	if err != nil {
		panic(err)
	}

	return db
}

type ReadWriteMyDBs struct {
	Read   *sqlx.DB
	Write  *sqlx.DB
	Delete *sqlx.DB
}

func MustNewMyDBs() ReadWriteMyDBs {
	return ReadWriteMyDBs{
		Read:   MustNewMySQL(config.MustMySQLReadConn()),
		Write:  MustNewMySQL(config.MustMySQLWriteConn()),
		Delete: MustNewMySQL(config.MustMySQLDeleteConn()),
	}
}
