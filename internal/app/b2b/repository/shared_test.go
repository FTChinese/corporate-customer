package repository

import (
	"github.com/FTChinese/ftacademy/internal/mock"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/enum"
	"testing"
)

func TestSharedRepo_SaveVersionedLicence(t *testing.T) {

	adm := mock.NewAdmin()
	p := mock.NewPersona()

	licToRenew := adm.StdLicenceBuilder().SetPersona(p).Build()

	r := NewSharedRepo(db.MockMySQL())

	type args struct {
		s licence.Versioned
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save versioned licence",
			args: args{
				s: adm.
					StdLicenceBuilder().
					Build().
					Versioned(licence.VersionActionCreate),
			},
			wantErr: false,
		},
		{
			name: "Save renewed licence",
			args: args{
				s: licToRenew.
					Renewed(price.MockPriceStdYear, pkg.TxnID()).
					Versioned(licence.VersionActionRenew).
					WithPriorVersion(licToRenew).
					WithMismatched(p.MemberBuilderFTC().
						WithPayMethod(enum.PayMethodApple).
						Build()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := r.SaveVersionedLicence(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SaveVersionedLicence() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
