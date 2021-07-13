package letter

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/postman"
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

func InvitationParcel(il licence.InvitedLicence, adminProfile admin.Profile) (postman.Parcel, error) {
	body, err := CtxInvitation{
		AssigneeName: il.Assignee.NormalizeName(),
		TeamName:     adminProfile.OrgName,
		Tier:         il.Plan.Tier,
		URL:          baseUrl + "/accept-invitation/" + il.Invitation.Token,
		AdminEmail:   adminProfile.Email,
	}.Render()

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   il.Assignee.Email.String,
		ToName:      il.Assignee.NormalizeName(),
		Subject:     "[FT中文网企业订阅]会员邀请",
		Body:        body,
	}, nil
}

// LicenceGrantedParcel write a letter to admin after
// a reader accepted an invitation and the corresponding
// licence is granted.
// We need to know the admin's account, reader's email
// the the licence's plan.
func LicenceGrantedParcel(il licence.InvitedLicence, adminAccount admin.BaseAccount) (postman.Parcel, error) {

	var data = CtxLicenceGranted{
		Name:           adminAccount.NormalizeName(),
		AssigneeEmail:  il.Assignee.Email.String,
		Tier:           il.Plan.Tier,
		ExpirationDate: il.Licence.ExpireDate,
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
