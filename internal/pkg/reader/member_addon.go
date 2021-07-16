package reader

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/dt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// AddOnInvoiceCreated is the result of manually adding an invoice.
type AddOnInvoiceCreated struct {
	Invoice    Invoice        `json:"invoice"`
	Membership Membership     `json:"membership"`
	Snapshot   MemberSnapshot `json:"snapshot"`
}

func (m Membership) HasAddOn() bool {
	return m.AddOn.Standard > 0 || m.AddOn.Premium > 0
}

func (m Membership) CarriedOverAddOn() addon.AddOn {
	return addon.New(m.Tier, m.RemainingDays())
}

// CarryOverInvoice creates a new invoice based on remaining days of current membership.
// This should only be used when user is upgrading from standard to premium using one-time purchase,
// or switch from one-time purchase to subscription mode.
func (m Membership) CarryOverInvoice() Invoice {
	return Invoice{
		ID:         pkg.InvoiceID(),
		CompoundID: m.CompoundID,
		Edition:    m.Edition,
		YearMonthDay: dt.YearMonthDay{
			Days: m.RemainingDays(),
		},
		AddOnSource:    addon.SourceCarryOver,
		AppleTxID:      null.String{},
		OrderID:        null.String{},
		OrderKind:      enum.OrderKindAddOn, // All carry-over invoice are add-ons
		PaidAmount:     0,
		PaymentMethod:  m.PaymentMethod,
		PriceID:        m.FtcPlanID,
		StripeSubsID:   null.String{},
		CreatedUTC:     chrono.TimeNow(),
		ConsumedUTC:    chrono.Time{}, // Will be consumed in the future.
		DateTimePeriod: dt.DateTimePeriod{},
		CarriedOverUtc: chrono.Time{},
	}
}
