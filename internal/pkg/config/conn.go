package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Conn represents a connection to a server or database.
type Conn struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

func GetConn(key string) (Conn, error) {
	var conn Conn
	err := viper.UnmarshalKey(key, &conn)
	if err != nil {
		return Conn{}, err
	}

	return conn, nil
}

// MustGetDBConn gets sql connection configuration.
func MustGetDBConn(c Config) Conn {
	var conn Conn
	var err error

	if c.Debug {
		conn, err = GetConn("mysql.dev")
	} else {
		conn, err = GetConn("mysql.master")
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using mysql server %s. Debugging: %t", conn.Host, c.Debug)

	return conn
}

// MustGEtEmailConn gets email server configuration.
func MustGetEmailConn() Conn {
	conn, err := GetConn("email.ftc")
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

// DSN build sql connection string.
func (c Conn) DSN() string {
	cfg := &mysql.Config{
		User:   c.User,
		Passwd: c.Pass,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		// Always use UTC time.
		// Pay attention to how string values are specified.
		// The string value provided to MySQL must be quoted in single quote for this driver to work,
		// which means the single quote itself must be included in the string value.
		// The resulting string value passed to MySQL should look like: `%27<you string value>%27`
		// See ASCII Encoding Reference https://www.w3schools.com/tags/ref_urlencode.asp
		Params: map[string]string{
			"time_zone": `'+00:00'`,
		},
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}

	return cfg.FormatDSN()
}

// NewDB creates opens a new sql connection.
func NewDB(c Conn) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", c.DSN())

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
	return db, nil
}
