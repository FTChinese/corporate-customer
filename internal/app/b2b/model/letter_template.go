package model

import (
	"text/template"
)

const letterTemplates = `
{{define "verification"}}
FT中文网企业订阅管理员 {{.Name}}，你好！

{{if .IsSignUp}}
感谢您注册FT中文网B2B订阅服务。您在此可以为您所属机构的成员订阅FT付费服务。请注意，B2B服务的账号独立于FT中文网的账号，如果您未使用此邮箱注册FT中文网账号，则此邮箱只能登录B2B服务。
{{end}}

请验证您的邮箱地址，帮助我们增强您的账号安全。

点击链接验证邮箱地址，如果链接无法点击，可以复制粘贴到浏览器地址栏：

{{.URL}}

您最近在FT中文网创建了新的B2B账号或更改了登录FT中文网B2B服务所用的邮箱，因此收到本邮件。如果您没有进行此操作，请忽略此邮件。

本邮件由系统自动生成，请勿回复。

FT中文网
{{end}}

{{define "passwordReset"}}
FT中文网B2B用户 {{.Name}}，你好！

获悉您遗失了B2B网站的登录密码，点击以下链接可以重置密码：

{{.URL}}

如果上述链接无法点击，可以复制粘贴到浏览器地址栏。

本链接3小时内有效。

本邮件由系统自动生成，请勿回复。

FT中文网
{{end}}

{{define "invitation"}}
FT中文网读者 {{.AssigneeName}}，你好！

{{.TeamName}}为您订阅了FT中文网会员 {.Tier.StringCN}}，请点击以下链接接受邀请。

{{.URL}}

接受邀请后即获得FT会员，可以阅读FT中文网的付费内容。

本链接3日内有效，请尽快接受邀请。如果链接已过期，请联系您所属机构的管理员 {{.AdminEmail}}。

本邮件由系统自动生成，请勿回复。

FT中文网
{{end}}

{{define "licenceGranted"}}
FT中文网B2B管理员 {{.Name}}，你好！

您通过FT中文网B2B业务邀请团队成员{{.AssigneeEmail}}成为FT中文网订阅用户，该成员已经接受了邀请，订阅方案的许可已经授予该用户：

订阅方案：{{.Tier.StringCN}}
到期日期：{{.ExpirationDate.String}}

您随时可以在B2B管理系统中撤销该用户的许可。

本邮件由系统自动生成，请勿回复。

FT中文网
{{end}}
`

var tmpl = template.Must(template.New("letter").Parse(letterTemplates))
