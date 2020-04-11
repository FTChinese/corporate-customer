package admin

import (
	"github.com/spf13/viper"
	"log"
)

var signingKey []byte

func init() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read viper config")
	}

	k := viper.GetString("web_app.b2b.jwt_signing_key")

	if k == "" {
		log.Fatal("JWT signing key not found")
	}

	signingKey = []byte(k)
}
