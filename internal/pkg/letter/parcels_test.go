package letter

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"testing"
)

func TestInvitationParcel(t *testing.T) {
	type args struct {
		assignee     licence.Assignee
		lic          licence.Licence
		adminProfile admin.Profile
	}
	tests := []struct {
		name    string
		args    args
		want    postman.Parcel
		wantErr bool
	}{
		{
			name: "Invitation parcel",
			args: args{
				assignee: licence.Assignee{
					FtcID:    null.StringFrom("86e4acb5-4884-4b7b-94d7-b9475b110996"),
					UnionID:  null.String{},
					Email:    null.StringFrom("neefrankie@163.com"),
					UserName: null.StringFrom("neefrankie"),
				},
				lic: licence.Licence{
					ID: "",
					Edition: price.Edition{
						Tier:  enum.TierStandard,
						Cycle: enum.CycleYear,
					},
					Creator:               admin.Creator{},
					Status:                0,
					CurrentPeriodStartUTC: chrono.Time{},
					CurrentPeriodEndUTC:   chrono.Time{},
					StartDateUTC:          chrono.Time{},
					TrialStartUTC:         chrono.Date{},
					TrialEndUTC:           chrono.Date{},
					HintGrantMismatch:     false,
					LatestTransactionID:   null.String{},
					LatestPrice:           price.Price{},
					LatestInvitation: licence.InvitationJSON{
						Invitation: licence.Invitation{
							ID:             "",
							Creator:        admin.Creator{},
							Status:         0,
							Description:    null.String{},
							ExpirationDays: 7,
							Email:          "",
							LicenceID:      "",
							Token:          ids.InvitationID(),
							RowTime:        admin.RowTime{},
						},
					},
					AssigneeID: null.String{},
					RowTime:    admin.RowTime{},
				},
				adminProfile: admin.Profile{
					BaseAccount: admin.BaseAccount{
						ID:          "",
						TeamID:      null.String{},
						Email:       "neefrankie@163.com",
						DisplayName: null.String{},
						Active:      false,
						Verified:    false,
					},
					TeamParams: input.TeamParams{
						OrgName:      "FTC",
						InvoiceTitle: null.String{},
						Phone:        null.String{},
					},
				},
			},
			want:    postman.Parcel{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InvitationParcel(tt.args.assignee, tt.args.lic, tt.args.adminProfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvitationParcel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("InvitationParcel() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%v", got)
		})
	}
}
