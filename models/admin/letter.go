package admin

// Letter is the data passed to template to generate the content
// of an email.
type Letter struct {
	URL      string // The link for verification, or password reset.
	IsSignUp bool   // Determines greeting.
}

const letterVerification = `
FT中文网B2B用户 {{.Name}}，你好！

{{if SignUp}}
感谢您注册FT中文网B2B订阅服务。您在此可以为您所属机构的成员订阅FT付费服务。请注意，B2B服务的账号独立于FT中文网的账号，如果您未使用此邮箱注册FT中文网账号，则此邮箱只能登录B2B服务。
{{end}}

请验证您的邮箱地址，帮助我们增强您的账号安全。

点击链接验证邮箱地址，如果链接无法点击，可以复制粘贴到浏览器地址栏：

{{.URL}}

您最近在FT中文网创建了新的B2B账号或更改了登录FT中文网B2B服务所用的邮箱，因此收到本邮件。如果您没有进行此操作，请忽略此邮件。

本邮件由系统自动生成，请勿回复。

FT中文网`

const letterPasswordReset = `
FT中文网B2B用户 {{.Name}}，你好！

获悉您遗失了B2B网站的登录密码，点击以下链接可以重置密码：

{{.URL}}

如果上述链接无法点击，可以复制粘贴到浏览器地址栏。

本链接3小时内有效。

FT中文网`
