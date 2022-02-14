package cmsrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_SavePayment(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		p checkout.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save payment",
			args: args{
				p: mock.OrderPaid().Payment,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SavePayment(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SavePayment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_SavePaymentOffer(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		offer checkout.PaymentOffer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save payment offer",
			args: args{
				offer: mock.OrderPaid().Offers[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SavePaymentOffer(tt.args.offer); (err != nil) != tt.wantErr {
				t.Errorf("SavePaymentOffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_SavePaymentResult(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		op checkout.OrderPaid
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save payment result",
			args: args{
				op: mock.OrderPaid(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SavePaymentResult(tt.args.op); (err != nil) != tt.wantErr {
				t.Errorf("SavePaymentResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_SavePaymentError(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		pe checkout.PaymentError
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save payment error",
			args: args{
				pe: checkout.PaymentError{
					TxnID:      ids.TxnID(),
					Message:    gofakeit.Sentence(20),
					CreatedUTC: chrono.TimeNow(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SavePaymentError(tt.args.pe); (err != nil) != tt.wantErr {
				t.Errorf("SavePaymentError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
