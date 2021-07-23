package letter

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/FTChinese/go-rest/chrono"
)

const fromAddress = "no-reply@ftchinese.com"

func VerificationParcel(ctx CtxVerification) (postman.Parcel, error) {
	body, err := ctx.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: fromAddress,
		FromName:    "FT中文网",
		ToAddress:   ctx.Email,
		ToName:      ctx.UserName,
		Subject:     "[FT中文网企业订阅]验证账号",
		Body:        body,
	}, nil
}

func PasswordResetParcel(a admin.BaseAccount, session admin.PwResetSession) (postman.Parcel, error) {
	body, err := CtxPwReset{
		UserName: a.NormalizeName(),
		Link:     session.BuildURL(),
		Duration: session.FormatDuration(),
	}.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: fromAddress,
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[FT中文网企业订阅]重置密码",
		Body:        body,
	}, nil
}

func InvitationParcel(assignee licence.Assignee, lic licence.BaseLicence, adminProfile admin.Profile) (postman.Parcel, error) {
	body, err := CtxInvitation{
		ToName:     assignee.NormalizeName(),
		Tier:       lic.Tier,
		URL:        pkg.B2BVerifyInvitationURL(lic.LatestInvitation.Token),
		AdminEmail: adminProfile.Email,
		TeamName:   adminProfile.OrgName,
	}.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   assignee.Email.String,
		ToName:      assignee.NormalizeName(),
		Subject:     "[FT中文网企业订阅]会员邀请",
		Body:        body,
	}, nil
}

// LicenceGrantedParcel write a letter to admin after
// a reader accepted an invitation and the corresponding
// licence is granted.
// We need to know the admin's account, reader's email
// the the licence's plan.
func LicenceGrantedParcel(lic licence.Licence, adminAccount admin.Profile) (postman.Parcel, error) {

	var data = CtxLicenceGranted{
		Name:           adminAccount.NormalizeName(),
		AssigneeEmail:  lic.Assignee.Email.String,
		Tier:           lic.Tier.StringCN(),
		ExpirationDate: chrono.DateFrom(lic.CurrentPeriodEndUTC.Time).String(),
	}

	body, err := data.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   adminAccount.Email,
		ToName:      data.Name,
		Subject:     "[FT中文网企业订阅]团队成员获得会员许可",
		Body:        body,
	}, nil
}
