package cmsrepo

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	gorest "github.com/FTChinese/go-rest"
	"go.uber.org/zap/zaptest"
	"testing"
)

func beforeTestOrder() checkout.OrderInputSchema {
	adm := mock.NewAdmin()

	schema := adm.CartBuilder().
		AddNewStandardN(5).
		AddRenewalStandardN(5).
		AddNewPremiumN(5).
		AddRenewalPremiumN(5).
		BuildOrderSchema()

	repo := mock.NewRepo()
	repo.InsertTeam(adm.Team())
	repo.InsertOrderSchema(schema)

	return schema
}

func TestEnv_listOrders(t *testing.T) {

	schema := beforeTestOrder()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		w    pkg.SQLWhere
		page gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []checkout.CMSOrderRow
		wantErr bool
	}{
		{
			name: "Retrieve all orders",
			args: args{
				w:    pkg.SQLWhere{},
				page: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
		{
			name: "Order of a team",
			args: args{
				w: checkout.NewOrderFilter(schema.OrderRow.TeamID).
					SQLWhere(),
				page: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.listOrders(tt.args.w, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("listOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("listOrders() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_countOrder(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		w pkg.SQLWhere
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Count orders",
			args: args{
				w: pkg.SQLWhere{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countOrder(tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("countOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("countOrder() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%d", got)
		})
	}
}

func TestEnv_ListOrders(t *testing.T) {

	schema := beforeTestOrder()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		filter checkout.OrderFilter
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    checkout.CMSOrderList
		wantErr bool
	}{
		{
			name: "List order",
			args: args{
				filter: checkout.NewOrderFilter(schema.OrderRow.TeamID),
				page:   gorest.NewPagination(1, 10),
			},
			want:    checkout.CMSOrderList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListOrders(tt.args.filter, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ListOrders() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_LoadOrder(t *testing.T) {

	schema := beforeTestOrder()

	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		args    args
		want    checkout.Order
		wantErr bool
	}{
		{
			name: "Load order",
			args: args{
				orderID: schema.OrderRow.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadOrder(tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadOrder() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
