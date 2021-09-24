package checkout

import (
	"github.com/FTChinese/go-rest/chrono"
	"sync"
)

const StmtSaveProcessingStats = `
INSERT INTO b2b.order_processing_log
SET order_id = :order_id,
	total_counter = :total_counter,
	success_counter = :success_counter,
	failure_counter = :failure_counter,
	start_utc = :start_utc,
	end_utc = :end_utc
`

type OrderProcessingStats struct {
	OrderID   string      `db:"order_id"`
	Total     int64       `db:"total_counter"`
	Succeeded int64       `db:"success_counter"`
	Failed    int64       `db:"failure_counter"`
	StartUTC  chrono.Time `db:"start_utc"`
	EndUTC    chrono.Time `db:"end_utc"`
	mux       sync.Mutex
}

func NewOrderProcessingStats(id string) *OrderProcessingStats {
	return &OrderProcessingStats{
		OrderID:  id,
		StartUTC: chrono.TimeNow(),
	}
}

func (l *OrderProcessingStats) IncTotal() {
	l.mux.Lock()
	l.Total++
	l.mux.Unlock()
}

func (l *OrderProcessingStats) IncSuccess() {
	l.mux.Lock()
	l.Succeeded++
	l.mux.Unlock()
}

func (l *OrderProcessingStats) IncFailure() {
	l.mux.Lock()
	l.Failed++
	l.mux.Unlock()
}
