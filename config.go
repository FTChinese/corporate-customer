package main

import (
	"github.com/FTChinese/b2b/database"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Debug   bool
	Version string
	BuiltAt string
	Year    int
}

func (c Config) MustGetDBConn(key string) database.Conn {
	var conn database.Conn
	var err error

	if c.Debug {
		err = viper.UnmarshalKey("mysql.dev", &conn)
	} else {
		err = viper.UnmarshalKey(key, &conn)
	}

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	logger.Infof("Using mysql server %s. Debugging: %t", conn.Host, c.Debug)

	return conn
}

func MustGetEmailConn() database.Conn {
	var emailConn database.Conn
	err := viper.UnmarshalKey("email.ftc", &emailConn)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	return emailConn
}

func MustGetSessionKey() string {
	k := viper.GetString("web_app.b2b.echo_session")

	if k == "" {
		logger.Error("Echo session key not found")
		os.Exit(1)
	}

	return k
}
