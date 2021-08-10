package letter

const (
	keyVrf               = "verification"
	keyVerified          = "email_verified"
	keyPwReset           = "password_reset"
	keyOrderCreated      = "order_created"
	keyLicenceInvitation = "licence_invitation"
	keyLicenceGranted    = "licence_granted"
)

const customerService = `
subscriber.service@ftchinese.com
`

var templates = map[string]string{
	keyVrf: `
FT中文网企业订阅管理员 {{.AdminName}}，你好！

感谢您注册FT中文网企业订阅服务。您可以为所属机构的批量订阅FT付费服务。

请验证您的邮箱地址，帮助我们增强您的账号安全。点击链接验证邮箱地址，如果链接无法点击，可以复制粘贴到浏览器地址栏：

{{.Link}}

您最近在FT中文网创建了新的B2B账号或更改了登录FT中文网B2B服务所用的邮箱，因此收到本邮件。如果您没有进行此操作，请忽略此邮件。

请注意，企业订阅服务的账号独立于FT中文网的账号，如果您未使用此邮箱注册FT中文网账号，则此邮箱只能登录企业订阅服务。

本邮件由系统自动生成，请勿回复。

FT中文网

-------------------------------
订阅咨询请联系：
` + customerService,

	keyPwReset: `
FT中文网企业订阅管理员 {{.AdminName}}，你好！

获悉您遗失了企业订阅网站的登录密码，点击以下链接可以重置密码：

{{.Link}}

如果上述链接无法点击，可以复制粘贴到浏览器地址栏。

本链接{{.Duration}}内有效。

本邮件由系统自动生成，请勿回复。

FT中文网`,
	keyOrderCreated: `
FT中文网企业订阅管理员 {{.AdminName}},

您的订单已创建：

{{.ID}}

{{range .Products}}
{{.Price.Tier | tierSC}}  {{.Price.UnitAmount | currency}}/{{.Price.Cycle.StringCN}}
	新增 {{.NewCopies}}份
	续订 {{.RenewalCopies}}份
{{end}}

共{{.ItemCount}}份，应付{{.AmountPayable | currency}}。

请联系我方客服洽谈付款事宜：
` + customerService,

	keyLicenceInvitation: `
FT中文网读者 {{.ReaderName}}，你好！

{{.TeamName}}为您订阅了FT中文网会员 {{.Tier}}，请点击以下链接接受邀请。

{{.Link}}

接受邀请后即获得FT会员，可以阅读FT中文网的付费内容。

本链接{{.Duration}}内有效，请尽快接受邀请。如果链接已过期，请联系您所属机构的管理员 {{.AdminEmail}}。

本邮件由系统自动生成，请勿回复。

FT中文网`,

	keyLicenceGranted: `
FT中文网B2B管理员 {{.AdminName}}，你好！

您通过FT中文网B2B业务邀请团队成员{{.AssigneeEmail}}成为FT中文网订阅用户，该成员已经接受了邀请，订阅方案的许可已经授予该用户：

订阅方案：{{.Tier}}
到期日期：{{.ExpirationDate}}

您随时可以在B2B管理系统中撤销该用户的许可。

本邮件由系统自动生成，请勿回复。

FT中文网`,
}
