package cmsrepo

import (
	"context"
	"database/sql"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"golang.org/x/sync/semaphore"
	"runtime"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
	sem        = semaphore.NewWeighted(int64(maxWorkers))
)

// listLicenceTxnPrice retrieves a list of licence transaction id and the price for this transaction,
// and send the resulting list to a channel.
// We do not retrieve the full transaction row here since
// each transaction needs to be locked upon confirmation.
func (env Env) listLicenceTxnPrice(orderID string) <-chan checkout.PriceOfLicenceTxn {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	ch := make(chan checkout.PriceOfLicenceTxn)

	go func() {
		defer close(ch)
		rows, err := env.DBs.Read.Queryx(checkout.StmtListTxnToConfirm, orderID)
		if err != nil {
			sugar.Error(err)
			return
		}

		var q checkout.PriceOfLicenceTxn
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
// Operations involved:
// * Lock licence transaction row;
// * Optionally lock renewal licence if the transaction is used to renew a licence;
// * Optionally lock membership is the licence-to-renew was granted to a membership;
// * Optionally retrieve (but not lock) assignee if the licence-to-renew is granted to someone;
// * Insert/update the newly created or renewed licence;
// * Mark transaction as finalized;
// * Optionally update membership for renewed licence;
// * Optionally save carry-over invoice if a one-time-purchase membership is overridden in renewal;
// * Save a versioned licence of both ante-change and post-change
// * Save a versioned membership of both ante-change and post-change.
func (env Env) buildLicence(txnPrice checkout.PriceOfLicenceTxn) (checkout.LicenceGenerated, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return checkout.LicenceGenerated{}, err
	}

	// Lock the row used to create/renew licence.
	licTxn, err := tx.LockLicenceTxn(txnPrice.TxnID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return checkout.LicenceGenerated{}, err
	}

	// If the transaction already finalized, stop.
	if licTxn.IsFinalized() {
		sugar.Infof("Transaction %s already finalized", licTxn.ID)
		_ = tx.Rollback()
		return checkout.LicenceGenerated{}, nil
	}

	sugar.Infof("Retrieved licence transaction %v", licTxn)

	// Retrieve current licence based on licence-to-renew.
	// Current licence is zero value for new licence.
	currLic, err := tx.LockLicenceCMS(licTxn.LicenceToRenew.ID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return checkout.LicenceGenerated{}, err
	}
	sugar.Infof("Retreived existing licence for transaction %s: %v", licTxn.ID, currLic)

	// Only exists if the licence is granted to a user
	// and the membership has not changed to auto-renewal
	// subscription.
	var curMmb reader.Membership
	var assignee licence.Assignee
	if currLic.IsGranted() {
		sugar.Infof("Licence %s is granted to %s", currLic.ID, currLic.AssigneeID.String)

		curMmb, err = tx.LockMember(currLic.AssigneeID.String)
		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return checkout.LicenceGenerated{}, err
		}
		sugar.Infof("Retrieve granted membership for licence %s: %v", currLic.ID, curMmb)

		// Retrieve assignee to complete licence queue data.
		// Here we should use the latest licence's assignee id
		// rather than the licence stored under LicenceToRenew
		// field since it might become obsolete as admin changes granting.
		assignee, err = env.RetrieveAssignee(currLic.AssigneeID.String)
		// Do not stop for error.
		if err != nil {
			// Fallback is assignee is not retrieved.
			if err == sql.ErrNoRows {
				assignee = licence.Assignee{
					FtcID:    curMmb.FtcID,
					UnionID:  curMmb.UnionID,
					Email:    null.String{},
					UserName: null.String{},
				}
			}
			sugar.Error(err)
		}
		sugar.Infof("Retrieved assignee %v for licence %s", assignee, currLic.ID)
	}

	// Build new licence,
	// updated licence queue.
	// optional new membership,
	// optional membership snapshot,
	// optional invoice
	result, err := checkout.GenerateLicence(checkout.LicenceGenParams{
		Price:     txnPrice.Price,
		LicTxn:    licTxn,
		CurLic:    currLic,
		Assignee:  assignee,
		CurMember: curMmb,
	})
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return checkout.LicenceGenerated{}, err
	}
	sugar.Infof("Licence generated %v", result.LicenceVersion.PostChange)
	sugar.Infof("Transactoin finalized %v", result.Transaction)
	sugar.Infof("Licence snapshot %v", result.LicenceVersion)
	sugar.Infof("AnteChange updated %v", result.MemberModified.MembershipVersion.PostChange.Membership)

	// Insert or update licence.
	sugar.Infof("Upserting licence %s", result.LicenceVersion.PostChange.ID)
	err = tx.UpsertLicence(licTxn.Kind, result.LicenceVersion.PostChange.Licence)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return checkout.LicenceGenerated{}, err
	}

	// This row in the queue table is finalized.
	sugar.Infof("Finalize licence transaction %s", result.Transaction.ID)
	err = tx.FinalizeLicenceTxn(result.Transaction)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return checkout.LicenceGenerated{}, err
	}

	// Optionally update membership
	if !result.MembershipVersion.PostChange.IsZero() {
		sugar.Infof("Updating membership %v", result.MemberModified.MembershipVersion.PostChange)
		err = tx.UpdateMember(result.MemberModified.MembershipVersion.PostChange.Membership)
		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return checkout.LicenceGenerated{}, err
		}
	}

	// Save optional invoice.
	if !result.CarryOverInvoice.IsZero() {
		sugar.Infof("Saving carried over invoice %v", result.CarryOverInvoice)

		err := tx.SaveInvoice(result.CarryOverInvoice)
		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return checkout.LicenceGenerated{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return checkout.LicenceGenerated{}, err
	}

	sugar.Infof("Save versioned licence %s", result.LicenceVersion.PostChange.ID)
	err = env.SaveVersionedLicence(result.LicenceVersion)
	if err != nil {
		sugar.Error(err)
	}

	// Version membership only if it is changed.
	if !result.MembershipVersion.IsZero() {
		sugar.Infof("Archiving membership %s", result.MemberModified.MembershipVersion.ID)
		err := env.ArchiveMembership(result.MemberModified.MembershipVersion)
		if err != nil {
			sugar.Error(err)
		}
	}

	return result, nil
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
	txnPriceCh := env.listLicenceTxnPrice(order.ID)

	queLog := checkout.NewOrderProcessingStats(order.ID)

	// Process each LicenceTransaction.
	for pq := range txnPriceCh {
		if err := sem.Acquire(ctx, 1); err != nil {
			sugar.Errorf("Failed to acquire semaphore: %v", err)
			break
		}

		sugar.Infof("------------------\nStart procession transaction %s", pq.TxnID)
		go func(p checkout.PriceOfLicenceTxn) {
			sugar.Infof("----- Start processing transaction %s -----", p.TxnID)

			queLog.IncTotal()
			_, err := env.buildLicence(p)
			if err != nil {
				sugar.Errorf("Error after processing %s, %s", p.TxnID, err)
				queLog.IncFailure()
				// Save the error.
				e := env.SavePaymentError(checkout.NewPaymentError(p.TxnID, err))
				if e != nil {
					sugar.Error(e)
				}
			} else {
				sugar.Infof("Transaction %s processed successfully", p.TxnID)
				queLog.IncSuccess()
				// TODO: send email to renewed membership.
			}
			sugar.Infof("----- Finished processing transaction %s -----", p.TxnID)
			sem.Release(1)
		}(pq)
	}

	// Acquire all the tokens to wait for any remaining workers to finish.
	//
	// If you are already waiting for the workers by some other means (such as an
	// errgroup.Group), you can omit this final Acquire call.
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		sugar.Infof("Failed to acquire semaphore: %v", err)
		return nil
	}

	sugar.Infof("Finished processing order %s", order.ID)

	queLog.EndUTC = chrono.TimeNow()

	sugar.Infof("Order queue finished %v", queLog)

	// Save stats of processing
	err := env.saveProcessingStats(queLog)
	if err != nil {
		sugar.Error(err)
	}

	// Update order status
	err = env.UpdateOrderStatus(order.ChangeStatus(checkout.StatusPaid))
	if err != nil {
		sugar.Error(err)
		return err
	}

	// TODO: send payment result to admin.
	return nil
}

// Save statistics of batching processing this order.
func (env Env) saveProcessingStats(l *checkout.OrderProcessingStats) error {
	_, err := env.DBs.Write.NamedExec(
		checkout.StmtSaveProcessingStats,
		l)
	if err != nil {
		return err
	}

	return nil
}
