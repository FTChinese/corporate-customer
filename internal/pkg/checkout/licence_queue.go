package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/ftacademy/pkg/sq"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

type LicenceQueue struct {
	ID           string         `json:"id" db:"id"`
	CreatedUTC   chrono.Time    `json:"createdUtc" db:"created_utc"`
	FinalizedUTC chrono.Time    `json:"finalizedUtc" db:"finalized_utc"`
	Index        int64          `json:"-" db:"array_index"`
	Kind         enum.OrderKind `json:"kind" db:"kind"`                  // Only create or renew.
	LicencePrior LicenceJSON    `json:"licencePrior" db:"licence_prior"` // The licence when this row is created. Do not use it to build renewed licence since it might already become obsolete.
	LicenceAfter LicenceJSON    `json:"licenceAfter" db:"licence_after"`
	OrderID      string         `json:"orderId" db:"order_id"`
	PriceID      string         `json:"priceId" db:"price_id"`
}

// NewLicenceQueue creates a new LicenceQueue for each
// copy newly created or renewed.
func NewLicenceQueue(orderID string, p price.Price, currLic licence.ExpandedLicence, i int) LicenceQueue {
	var k enum.OrderKind
	if currLic.ID == "" {
		k = enum.OrderKindCreate
	} else {
		k = enum.OrderKindRenew
	}

	return LicenceQueue{
		ID:           pkg.LicenceQueueID(),
		FinalizedUTC: chrono.Time{},
		Index:        int64(i),
		Kind:         k,
		LicencePrior: LicenceJSON{currLic},
		LicenceAfter: LicenceJSON{},
		OrderID:      orderID,
		PriceID:      p.ID,
	}
}

// Finalize after a licence is created or renewed.
func (q LicenceQueue) Finalize(lic licence.ExpandedLicence) LicenceQueue {
	q.FinalizedUTC = chrono.TimeNow()
	q.LicenceAfter = LicenceJSON{lic}

	return q
}

// RowValues build the values of an SQL bulk insert.
func (q LicenceQueue) RowValues() []interface{} {
	return []interface{}{
		q.ID,
		q.FinalizedUTC,
		q.Index,
		q.Kind,
		q.LicencePrior,
		q.LicenceAfter,
		q.OrderID,
		q.PriceID,
	}
}

// BulkLicenceQueue is used to build the values part of
// bulk insert, together with StmtBulkLicenceQueue.
type BulkLicenceQueue []LicenceQueue

func (b BulkLicenceQueue) Each(handler func(row sq.BulkInsertRow)) {
	for _, v := range b {
		handler(v)
	}
}

// GroupedQueues groups LicenceQueue for the specified
// price group of an order.
type GroupedQueues struct {
	PriceID  string         `json:"priceId"`
	Creation []LicenceQueue `json:"creation"`
	Renewal  []LicenceQueue `json:"renewal"`
}

func NewGroupedQueues(priceID string, rows []LicenceQueue) GroupedQueues {
	var g = GroupedQueues{
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

// PriceOfQueue maps the LicenceQueue id to the price it used.
type PriceOfQueue struct {
	QueueID string      `db:"queue_id"`
	Price   price.Price `db:"price"`
}
