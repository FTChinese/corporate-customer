package config

import (
	"errors"
	"github.com/spf13/viper"
)

type WechatApp struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"secret"`
}

func (a WechatApp) Validate() error {
	if a.AppID == "" || a.AppSecret == "" {
		return errors.New("wechat oauth app id or secret cannot be empty")
	}

	return nil
}

func LoadWechatApp(key string) (WechatApp, error) {
	var app WechatApp
	err := viper.UnmarshalKey(key, &app)
	if err != nil {
		return WechatApp{}, err
	}

	if err := app.Validate(); err != nil {
		return WechatApp{}, err
	}

	return app, nil
}

func MustLoadWechatApp(key string) WechatApp {
	app, err := LoadWechatApp(key)
	if err != nil {
		panic(err)
	}

	return app
}

func MustWxWebApp() WechatApp {
	return MustLoadWechatApp("wxapp.web_oauth")
}
