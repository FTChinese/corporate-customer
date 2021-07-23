// +build !production

package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

func MockAssignee() Assignee {
	faker.SeedGoFake()

	return Assignee{
		FtcID:    null.StringFrom(uuid.New().String()),
		UnionID:  null.String{},
		Email:    null.StringFrom(gofakeit.Email()),
		UserName: null.StringFrom(gofakeit.Username()),
	}
}

func MockLicence(p price.Price, creator admin.PassportClaims) Licence {
	return Licence{
		BaseLicence: NewBaseLicence(p, pkg.OrderID(), creator),
		Assignee:    MockAssignee(),
	}
}
