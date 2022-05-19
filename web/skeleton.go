package web

var siteURL = "https://www.ftacademy.cn"

type HyperLink struct {
	Href string
	Name string
}

type FooterColumn struct {
	Title string
	Items []HyperLink
}

type Footer struct {
	Matrix        []FooterColumn
	Year          int
	ClientVersion string
	ServerVersion string
}

type Head struct {
	SiteURL     string
	LogoURL     string
	IconBaseURL string
	IconSizes   []string
	StyleLinks  []string
}

var DefaultHead = Head{
	SiteURL:     siteURL,
	LogoURL:     siteURL + "/images/logo/brand-ftc-masthead.svg",
	IconBaseURL: siteURL + "/images/favicons",
	IconSizes: []string{
		"180x180",
		"152x152",
		"120x120",
		"76x76",
	},
	StyleLinks: []string{
		"https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/5.1.3/css/bootstrap.min.css",
		"https://www.ftacademy.cn/frontend/styles/ftc-bootstrap.css",
	},
}

var DefaultFooterCols = []FooterColumn{
	{
		Title: "支持",
		Items: []HyperLink{
			{
				Name: "关于我们",
				Href: "",
			},
			{
				Name: "职业机会",
				Href: "https://www.ftchinese.com/jobs/?from=ft",
			},
			{
				Name: "问题回馈",
				Href: "https://www.ftchinese.com/m/corp/faq.html",
			},
			{
				Name: "联系方式",
				Href: "https://www.ftchinese.com/m/corp/contact.html",
			},
		},
	},
	{
		Title: "法律事务",
		Items: []HyperLink{
			{
				Name: "服务条款",
				Href: "https://www.ftchinese.com/m/corp/service.html",
			},
			{
				Name: "版权声明",
				Href: "https://www.ftchinese.com/m/corp/copyright.html",
			},
		},
	},
	{
		Title: "服务",
		Items: []HyperLink{
			{
				Name: "企业订阅",
				Href: "https://next.ftacademy.cn/corporate/login",
			},
			{
				Name: "广告业务",
				Href: "https://www.ftchinese.com/m/corp/sales.html",
			},
			{
				Name: "会议活动",
				Href: "https://www.ftchinese.com/m/events/event.html",
			},
			{
				Name: "会员信息中心",
				Href: "https://next.ftacademy.cn/reader/login",
			},
			{
				Name: "最新动态",
				Href: "https://www.ftchinese.com/m/marketing/ftc.html",
			},
			{
				Name: "合作伙伴",
				Href: "https://www.ftchinese.com/m/corp/partner.html",
			},
		},
	},
	{
		Title: "关注我们",
		Items: []HyperLink{
			{
				Name: "微信",
				Href: "https://www.ftchinese.com/m/corp/follow.html",
			},
			{
				Name: "微博",
				Href: "https://weibo.com/ftchinese",
			},
			{
				Name: "Linkedin",
				Href: "https://www.linkedin.com/company/4865254?trk=hp-feed-company-Name",
			},
			{
				Name: "Facebook",
				Href: "https://www.facebook.com/financialtimeschinese",
			},
			{
				Name: "Twitter",
				Href: "https://twitter.com/FTChinese",
			},
		},
	},
	{
		Title: "FT产品",
		Items: []HyperLink{
			{
				Name: "FT研究院",
				Href: "https://www.ftchinese.com/m/marketing/intelligence.html",
			},
			{
				Name: "FT商学院",
				Href: "https://www.ftchinese.com/channel/mba.html",
			},
			{
				Name: "FT电子书",
				Href: "https://www.ftchinese.com/m/marketing/ebook.html",
			},
			{
				Name: "数据新闻",
				Href: "https://www.ftchinese.com/channel/datanews.html",
			},
			{
				Name: "FT英文版",
				Href: "https://www.ft.com/",
			},
		},
	},
	{
		Title: "移动应用",
		Items: []HyperLink{
			{
				Name: "安卓",
				Href: "https://next.ftchinese.com/android/latest",
			},
		},
	},
}
