package checkout

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/ftacademy/pkg/sq"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

type LicenceTransaction struct {
	ID             string              `json:"id" db:"txn_id"`
	Kind           enum.OrderKind      `json:"kind" db:"kind"`                       // Only create or renew.
	LicenceToRenew ExpandedLicenceJSON `json:"licenceToRenew" db:"licence_to_renew"` // The licence when this row is created. Do not use it to build renewed licence since it might already become obsolete.
	OrderID        string              `json:"orderId" db:"order_id"`
	PriceID        string              `json:"priceId" db:"price_id"`
	admin.Creator
	CreatedUTC   chrono.Time `json:"createdUtc" db:"created_utc"`
	FinalizedUTC chrono.Time `json:"finalizedUtc" db:"finalized_utc"`
}

// NewLicenceTransaction creates a new LicenceTransaction for each
// copy newly created or renewed.
func NewLicenceTransaction(
	orderID string,
	p price.Price,
	by admin.Creator,
	currLic licence.ExpandedLicence,
) LicenceTransaction {
	var k enum.OrderKind
	if currLic.ID == "" {
		k = enum.OrderKindCreate
	} else {
		k = enum.OrderKindRenew
	}

	return LicenceTransaction{
		ID:             ids.TxnID(),
		Kind:           k,
		LicenceToRenew: ExpandedLicenceJSON{currLic},
		OrderID:        orderID,
		PriceID:        p.ID,
		Creator:        by,
		CreatedUTC:     chrono.TimeNow(),
		FinalizedUTC:   chrono.Time{},
	}
}

func (t LicenceTransaction) BuildLicence(current licence.Licence, p price.Price) (licence.Licence, error) {
	switch t.Kind {
	case enum.OrderKindCreate:
		return licence.NewLicence(p, t.ID, t.Creator), nil

	case enum.OrderKindRenew:
		return current.Renewed(p, t.ID), nil
	}

	return licence.Licence{}, errors.New("unknown order kind")
}

func (t LicenceTransaction) IsFinalized() bool {
	return !t.FinalizedUTC.IsZero()
}

// Finalize after a licence is created or renewed.
func (t LicenceTransaction) Finalize() LicenceTransaction {
	t.FinalizedUTC = chrono.TimeNow()

	return t
}

// RowValues build the values of an SQL bulk insert.
func (t LicenceTransaction) RowValues() []interface{} {
	return []interface{}{
		t.ID,
		t.Kind,
		t.LicenceToRenew,
		t.OrderID,
		t.PriceID,
		t.AdminID,
		t.TeamID,
		t.CreatedUTC,
		t.FinalizedUTC,
	}
}

// BulkLicenceTxn is used to build the values part of
// bulk insert, together with StmtBulkLicenceTxn.
type BulkLicenceTxn []LicenceTransaction

func (b BulkLicenceTxn) Each(handler func(row sq.BulkInsertRow)) {
	for _, v := range b {
		handler(v)
	}
}

// GroupedLicenceTxn groups LicenceTransaction for the specified
// price group of an order.
type GroupedLicenceTxn struct {
	PriceID  string               `json:"priceId"`
	Creation []LicenceTransaction `json:"creation"`
	Renewal  []LicenceTransaction `json:"renewal"`
}

func NewGroupedTxn(priceID string, rows []LicenceTransaction) GroupedLicenceTxn {
	var g = GroupedLicenceTxn{
		PriceID: priceID,
	}

	for _, item := range rows {
		if item.PriceID != priceID {
			continue
		}

		switch item.Kind {
		case enum.OrderKindCreate:
			g.Creation = append(g.Creation, item)
		case enum.OrderKindRenew:
			g.Renewal = append(g.Renewal, item)
		}
	}

	return g
}
