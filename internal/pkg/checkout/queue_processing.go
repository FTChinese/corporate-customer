package checkout

import (
	"github.com/FTChinese/go-rest/chrono"
	"sync"
)

const StmtSaveProcessingLog = `
INSERT INTO b2b.queue_processing_log
SET total_counter = :total_counter,
	success_counter = :success_counter,
	failure_counter = :failure_counter,
	start_utc = :start_utc,
	end_utc = :end_utc
`

type QueueProcessingLog struct {
	Total     int64       `db:"total_counter"`
	Succeeded int64       `db:"success_counter"`
	Failed    int64       `db:"failure_counter"`
	StartUTC  chrono.Time `db:"start_utc"`
	EndUTC    chrono.Time `db:"end_utc"`
	mux       sync.Mutex
}

func NewQueueProcessingLog() *QueueProcessingLog {
	return &QueueProcessingLog{
		StartUTC: chrono.TimeNow(),
	}
}

func (l *QueueProcessingLog) IncTotal() {
	l.mux.Lock()
	l.Total++
	l.mux.Unlock()
}

func (p *QueueProcessingLog) IncSuccess() {
	p.mux.Lock()
	p.Succeeded++
	p.mux.Unlock()
}

func (p *QueueProcessingLog) IncFailure() {
	p.mux.Lock()
	p.Failed++
	p.mux.Unlock()
}
