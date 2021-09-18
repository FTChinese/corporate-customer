package faker

import "testing"

func TestGenStripeSubID(t *testing.T) {
	t.Logf("Stripe subscription id %s", GenStripeSubID())
}

func TestGenStripeCusID(t *testing.T) {
	t.Logf("Stripe customer id %s", GenStripeCusID())
}

func TestGenAppleSubID(t *testing.T) {
	t.Logf("IAP tx id %s", GenAppleSubID())
}
