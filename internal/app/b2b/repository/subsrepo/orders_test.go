package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/faker"
	gorest "github.com/FTChinese/go-rest"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestEnv_CreateOrder(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		schema checkout.OrderInputSchema
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create order",
			args: args{
				schema: checkout.MockOrderInputSchema(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreateOrder(tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_listOrders(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	s := checkout.MockOrderInputSchema()

	_ = env.CreateOrder(s)

	type args struct {
		teamID string
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []checkout.OrderRow
		wantErr bool
	}{
		{
			name: "List orders",
			args: args{
				teamID: s.OrderRow.TeamID,
				page:   gorest.NewPagination(1, 20),
			},
			want: []checkout.OrderRow{
				s.OrderRow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.listOrders(tt.args.teamID, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("listOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_countOrder(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	s := checkout.MockOrderInputSchema()

	_ = env.CreateOrder(s)

	type args struct {
		teamID string
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
				teamID: s.OrderRow.TeamID,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countOrder(tt.args.teamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("countOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_ListOrders(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	s := checkout.MockOrderInputSchema()

	_ = env.CreateOrder(s)

	type args struct {
		teamID string
		page   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    checkout.OrderRowList
		wantErr bool
	}{
		{
			name: "List orders",
			args: args{
				teamID: s.OrderRow.TeamID,
				page:   gorest.NewPagination(1, 10),
			},
			want: checkout.OrderRowList{
				PagedList: pkg.PagedList{
					Total:      1,
					Pagination: gorest.NewPagination(1, 10),
					Err:        nil,
				},
				Data: []checkout.OrderRow{
					s.OrderRow,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListOrders(tt.args.teamID, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_orderDetails(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	s := checkout.MockOrderInputSchema()

	_ = env.CreateOrder(s)

	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		args    args
		want    checkout.Order
		wantErr bool
	}{
		{
			name: "Retrieve an order",
			args: args{
				r: admin.AccessRight{
					RowID:  s.OrderRow.ID,
					TeamID: s.OrderRow.TeamID,
				},
			},
			want: checkout.Order{
				BaseOrder: s.OrderRow.BaseOrder,
				Items:     nil,
				Payment:   checkout.Payment{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.orderDetails(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderDetails() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_orderItems(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	s := checkout.MockOrderInputSchema()

	_ = env.CreateOrder(s)

	t.Logf("%+v", s.ItemRows)

	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		args    args
		want    []checkout.OrderItem
		wantErr bool
	}{
		{
			name: "Retrieve order items",
			args: args{
				orderID: s.OrderRow.ID,
			},
			want:    s.ItemRows,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.orderItems(tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("orderItems() got = %v, \nwant %v", got, tt.want)
			//}
			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_LoadOrderDetails(t *testing.T) {
	env := NewEnv(db.MockMySQL(), zaptest.NewLogger(t))

	s := checkout.MockOrderInputSchema()

	_ = env.CreateOrder(s)

	type args struct {
		r admin.AccessRight
	}
	tests := []struct {
		name    string
		args    args
		want    checkout.Order
		wantErr bool
	}{
		{
			name: "Order details",
			args: args{
				r: admin.AccessRight{
					RowID:  s.OrderRow.ID,
					TeamID: s.OrderRow.TeamID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadDetailedOrder(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDetailedOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadDetailedOrder() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
