package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

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
	log.Printf("Pick %s for %s", a.name, prodDev[prod])

	if prod {
		return a.Prod
	}

	return a.Dev
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
	log.Printf("Loading API access token...")
	return MustLoadAPIConfig("api_keys.ftacademy")
}

func MustProdAPIv6BaseURL() API {
	log.Printf("Loading production API v6 base url...")
	return MustLoadAPIConfig("api_urls.subs_v6")
}

func MustSandboxAPIURL() API {
	log.Printf("Loading sandbox API base url...")
	return MustLoadAPIConfig("api_urls.sandbox")
}

func MustStripePubKey() API {
	log.Printf("Loading stripe publishable key...")
	return MustLoadAPIConfig("api_keys.stripe_publishable")
}
