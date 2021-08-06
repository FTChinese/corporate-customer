package config

import (
	"bytes"
	"github.com/spf13/viper"
)

func SetupViper(b []byte) error {
	viper.SetConfigType("toml")
	//viper.SetConfigName("api")
	//viper.AddConfigPath("$HOME/config")

	err := viper.ReadConfig(bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	return nil
}

func MustSetupViper(b []byte) {
	if err := SetupViper(b); err != nil {
		panic(err)
	}
}
