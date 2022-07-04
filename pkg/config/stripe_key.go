package config

import (
	"log"
)

type StripeKeyHolder struct {
	Key  string `json:"key"`
	Live bool   `json:"live"`
}

type StripeKeyStore struct {
	Sandbox StripeKeyHolder
	Live    StripeKeyHolder
}

func NewStripePubKeys() StripeKeyStore {
	keys := MustStripePubKey()
	return StripeKeyStore{
		Sandbox: StripeKeyHolder{
			Key:  keys.Dev,
			Live: false,
		},
		Live: StripeKeyHolder{
			Key:  keys.Prod,
			Live: true,
		},
	}
}

func (s StripeKeyStore) Select(live bool) StripeKeyHolder {
	log.Printf("Selecting stripe publishable key for %s\n", liveTest[live])
	if live {
		return s.Live
	}

	return s.Sandbox
}
