package letter

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

const (
	fromAddress = "no-reply@ftchinese.com"
	fromName    = "FT中文网"
	subjectName = "[FT中文网企业订阅]"
)

func VerificationParcel(a admin.BaseAccount, v admin.EmailVerifier) (postman.Parcel, error) {
	ctx := CtxVerification{
		Email:     a.Email,
		AdminName: a.NormalizeName(),
		Link:      v.BuildURL(),
	}

	body, err := ctx.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: fromAddress,
		FromName:    fromName,
		ToAddress:   ctx.Email,
		ToName:      ctx.AdminName,
		Subject:     subjectName + "验证账号",
		Body:        body,
	}, nil
}

func PasswordResetParcel(a admin.BaseAccount, session admin.PwResetSession) (postman.Parcel, error) {
	body, err := CtxPwReset{
		AdminName: a.NormalizeName(),
		Link:      session.BuildURL(),
		Duration:  session.FormatDuration(),
	}.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: fromAddress,
		FromName:    fromName,
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     subjectName + "重置密码",
		Body:        body,
	}, nil
}

func OrderCreatedParcel(a admin.Profile, order checkout.OrderRow) (postman.Parcel, error) {
	name := a.NormalizeName()

	body, err := CtxOrderCreated{
		AdminName: name,
		OrderRow:  order,
	}.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: fromAddress,
		FromName:    fromName,
		ToAddress:   a.Email,
		ToName:      name,
		Subject:     subjectName + "订单",
		Body:        body,
	}, nil
}

func InvitationParcel(assignee licence.Assignee, lic licence.BaseLicence, adminProfile admin.Profile) (postman.Parcel, error) {
	// If assignee does not exist.
	if assignee.IsZero() {
		assignee.Email = null.StringFrom(lic.LatestInvitation.Email)
	}

	body, err := CtxInvitation{
		ReaderName: assignee.NormalizeName(),
		AdminEmail: adminProfile.Email,
		TeamName:   adminProfile.OrgName,
		Tier:       lic.Tier.StringCN(),
		Link:       pkg.B2BVerifyInvitationURL(lic.LatestInvitation.Token),
		Duration:   lic.LatestInvitation.FormatDuration(),
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
		AdminName:      adminAccount.NormalizeName(),
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
		ToName:      data.AdminName,
		Subject:     "[FT中文网企业订阅]团队成员获得会员许可",
		Body:        body,
	}, nil
}
