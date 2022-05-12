package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type KeyHolder struct {
	Live string
	Test string
}

func (h KeyHolder) Select(live bool) string {
	if live {
		return h.Live
	}

	return h.Test
}

// API is used to hold api access token or api base url.
// Those keys are always comes in pair, one for development and one for production.
type API struct {
	Dev  string `mapstructure:"dev"`
	Prod string `mapstructure:"prod"`
	name string
}

func (a API) Validate() error {
	if a.Dev == "" || a.Prod == "" {
		return errors.New("dev or prod key not found")
	}

	return nil
}

func (a API) Pick(prod bool) string {
	log.Printf("Using %s for production %t", a.name, prod)

	if prod {
		return a.Prod
	}

	return a.Dev
}

func (a API) KeyHolder(prod bool) KeyHolder {
	// Production server provides live/test versions.
	// Live version uses production key;
	// test version uses development key.
	if prod {
		return KeyHolder{
			Live: a.Prod,
			Test: a.Dev,
		}
	}

	// If the binary is running under development environment,
	// live/test both use development key.
	return KeyHolder{
		Live: a.Dev,
		Test: a.Dev,
	}
}

func LoadAPIConfig(name string) (API, error) {
	var keys API
	err := viper.UnmarshalKey(name, &keys)
	if err != nil {
		return keys, err
	}

	if err := keys.Validate(); err != nil {
		return keys, err
	}

	keys.name = name

	return keys, nil
}

func MustLoadAPIConfig(name string) API {
	k, err := LoadAPIConfig(name)
	if err != nil {
		log.Fatalf("cannot get %s: %s", name, err.Error())
	}

	return k
}

func MustSubsAPIKey() API {
	return MustLoadAPIConfig("api_keys.ftacademy")
}

func MustSubsAPIv6BaseURL() API {
	return MustLoadAPIConfig("api_urls.subs_v6")
}

func MustAPISandboxURL() API {
	return MustLoadAPIConfig("api_urls.sandbox")
}

func MustStripePubKey() API {
	return MustLoadAPIConfig("api_keys.stripe_publishable")
}
