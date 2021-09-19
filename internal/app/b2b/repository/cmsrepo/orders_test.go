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

func TestEnv_listOrders(t *testing.T) {
	adm := mock.NewAdmin()

	repo := mock.NewRepo()
	repo.InsertTeam(adm.Team())

	schema := adm.CartBuilder().
		AddNewStandardN(5).
		AddRenewalStandardN(5).
		AddNewPremiumN(5).
		AddRenewalPremiumN(5).
		BuildOrderSchema()

	repo.InsertOrderSchema(schema)

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
				w: checkout.NewOrderFilter(adm.TeamID.String).
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
