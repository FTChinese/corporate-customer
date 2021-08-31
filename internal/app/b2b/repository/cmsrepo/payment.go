package cmsrepo

import "github.com/FTChinese/ftacademy/internal/pkg/checkout"

func (env Env) SavePayment(p checkout.Payment) error {
	_, err := env.DBs.Write.NamedExec(
		checkout.StmtSavePayment,
		p)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) SavePaymentOffer(offer checkout.PaymentOffer) error {
	_, err := env.DBs.Write.NamedExec(
		checkout.StmtSavePaymentOffer,
		offer)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) SavePaymentResult(op checkout.OrderPaid) error {
	err := env.SavePayment(op.Payment)
	if err != nil {
		return err
	}

	for _, v := range op.Offers {
		err := env.SavePaymentOffer(v)
		if err != nil {
			return err
		}
	}

	return nil
}
