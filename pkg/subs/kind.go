package subs

type Kind int

const (
	KindZero Kind = iota
	KindRetailNew
	KindRetailRenew
	KindRetailUpgrade   // Switching subscription tier, e.g., from standard to premium.
	KindRetailAddOn     // User purchased addon
	KindOverrideOneTime // One-time purchase could be overridden by b2b or auto-renewal.
	KindB2BNew
	KindB2BRenew
	KindB2BSwitchLicence
	KindAutoRenewSwitchCycle // Switching subscription billing cycle, e.g., from month to year.
)
