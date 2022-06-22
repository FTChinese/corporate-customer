package config

import (
	"github.com/spf13/viper"
	"log"
)

// Connect represents a connection to a server or database.
type Connect struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

func GetConn(key string) (Connect, error) {
	var conn Connect
	err := viper.UnmarshalKey(key, &conn)
	if err != nil {
		return Connect{}, err
	}

	return conn, nil
}

func MustMySQLConn(key string, prod bool) Connect {
	var conn Connect
	var err error

	if prod {
		conn, err = GetConn(key)
	} else {
		conn, err = GetConn("mysql.dev")
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got mysql ip %s", conn.Host)

	return conn
}

func MustMySQLReadConn(prod bool) Connect {
	log.Printf("Initializting MySQL read DB in %s", prodDev[prod])
	return MustMySQLConn("mysql.read", prod)
}

func MustMySQLWriteConn(prod bool) Connect {
	log.Printf("Intializng MySQL write DB in %s", prodDev[prod])
	return MustMySQLConn("mysql.write", prod)
}

func MustMySQLDeleteConn(prod bool) Connect {
	log.Printf("Intializng MySQL delete DB in %s", prodDev[prod])
	return MustMySQLConn("mysql.delete", prod)
}

func MustGetHanqiConn() Connect {
	var conn Connect
	err := viper.UnmarshalKey("email.hanqi", &conn)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
