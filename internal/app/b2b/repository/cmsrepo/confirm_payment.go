package cmsrepo

import (
	"context"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/go-rest/chrono"
	"golang.org/x/sync/semaphore"
	"runtime"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
	sem        = semaphore.NewWeighted(int64(maxWorkers))
)

func (env Env) retrieveOrderQueue(orderID string) <-chan checkout.PriceOfQueue {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	ch := make(chan checkout.PriceOfQueue)

	go func() {
		defer close(ch)
		rows, err := env.DBs.Read.Queryx(checkout.StmtListPriceOfQueue, orderID)
		if err != nil {
			sugar.Error(err)
			return
		}

		var q checkout.PriceOfQueue
		for rows.Next() {
			err := rows.StructScan(&q)
			if err != nil {
				sugar.Error(err)
				continue
			}

			ch <- q
		}
	}()

	return ch
}

// buildLicence creates/updates a licence based on
// a row from licence queue.
// After a licence is created/renewed, we should do:
// * Tell licence user membership is renewed if it exists.
func (env Env) buildLicence(order checkout.Order, pq checkout.PriceOfQueue) error {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return err
	}

	// Lock the row used to create/renew licence.
	licQueue, err := tx.LockLicenceQueue(pq.QueueID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Retrieve current licence based on the queue's prior
	// field. Current licence is zero value for new
	// licence.
	currLic, err := tx.LockBaseLicenceCMS(licQueue.LicencePrior.ID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Collect data and build new licence.
	newLic, err := checkout.LicenceBuilder{
		Kind:           licQueue.Kind,
		CurrentLicence: currLic,
		Price:          pq.Price,
		Order:          order,
	}.Build()
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Insert or update licence.
	err = tx.UpsertLicence(licQueue.Kind, newLic)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return err
	}

	// If current licence is already granted to a user,
	// update membership.
	renewed, err := tx.RenewMembership(newLic)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
	}

	// Retrieve assignee to complete licence queue data.
	// Here we should use the latest licence's assignee id
	// rather than the licence stored under LicencePrior
	// field since it might become obsolete as admin changes
	// granting.
	var assignee licence.Assignee
	if newLic.AssigneeID.Valid {
		assignee, err = env.RetrieveAssignee(newLic.AssigneeID.String)
		if err != nil {
			sugar.Error(err)
		}
	}

	// This row in the queue table is finalized.
	err = tx.FinalizeLicenceQueue(licQueue.Finalize(licence.ExpandedLicence{
		Licence:  newLic,
		Assignee: licence.AssigneeJSON{Assignee: assignee},
	}))

	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Archive previous membership.
	if !renewed.Archive.Membership.IsZero() {
		err := env.ArchiveMembership(renewed.Archive)
		if err != nil {
			sugar.Error(err)
		}
	}

	return nil
}

// ConfirmPayment creates/renew licences under an order
// one by one.
// After the job is finished:
// * Tell admin the job is done;
// * Tell cms the job is done;
func (env Env) ConfirmPayment(order checkout.Order) error {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()
	ctx := context.Background()

	// Retrieve queued licence of an order.
	pqCh := env.retrieveOrderQueue(order.ID)

	queLog := checkout.NewQueueProcessingLog()

	// Process each LicenceQueue.
	for pq := range pqCh {
		if err := sem.Acquire(ctx, 1); err != nil {
			sugar.Errorf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(o checkout.Order, p checkout.PriceOfQueue) {
			queLog.IncTotal()
			err := env.buildLicence(o, p)
			if err != nil {
				sugar.Error(err)
				queLog.IncFailure()
			} else {
				queLog.IncSuccess()
			}
		}(order, pq)
	}

	// Acquire all the tokens to wait for any remaining workers to finish.
	//
	// If you are already waiting for the workers by some other means (such as an
	// errgroup.Group), you can omit this final Acquire call.
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		sugar.Infof("Failed to acquire semaphore: %v", err)
		return nil
	}

	queLog.EndUTC = chrono.TimeNow()

	err := env.saveProcessingLog(queLog)
	if err != nil {
		sugar.Error(err)
	}

	sugar.Infof("Order queue finished %v", queLog)

	// Update order status
	updatedOrder := order.ChangeStatus(checkout.StatusPaid)
	err = env.UpdateOrderStatus(updatedOrder)
	if err != nil {
		sugar.Error(err)
		return err
	}

	return nil
}

func (env Env) saveProcessingLog(l *checkout.QueueProcessingLog) error {
	_, err := env.DBs.Write.NamedExec(
		checkout.StmtSaveProcessingLog,
		l)
	if err != nil {
		return err
	}

	return nil
}
