package views

import "github.com/FTChinese/b2b/widget"

type FooterColumn struct {
	Title string
	Items []widget.Link
}

var Footer = []FooterColumn{
	{
		Title: "支持",
		Items: []widget.Link{
			{
				Text: "关于我们",
				Href: "http://www.ftchinese.com/m/corp/aboutus.html",
			},
			{
				Text: "职业机会",
				Href: "http://www.ftchinese.com/jobs/?from=ft",
			},
			{
				Text: "问题回馈",
				Href: "http://www.ftchinese.com/m/corp/faq.html",
			},
			{
				Text: "联系方式",
				Href: "http://www.ftchinese.com/m/corp/contact.html",
			},
		},
	},
	{
		Title: "法律事务",
		Items: []widget.Link{
			{
				Text: "服务条款",
				Href: "http://www.ftchinese.com/m/corp/service.html",
			},
			{
				Text: "版权声明",
				Href: "http://www.ftchinese.com/m/corp/copyright.html",
			},
		},
	},
	{
		Title: "服务",
		Items: []widget.Link{
			{
				Text: "广告业务",
				Href: "http://www.ftchinese.com/m/corp/sales.html",
			},
			{
				Text: "会议活动",
				Href: "http://www.ftchinese.com/m/events/event.html",
			},
			{
				Text: "会员信息中心",
				Href: "http://www.ftchinese.com/m/marketing/home.html",
			},
			{
				Text: "最新动态",
				Href: "http://www.ftchinese.com/m/marketing/ftc.html",
			},
			{
				Text: "合作伙伴",
				Href: "http://www.ftchinese.com/m/corp/partner.html",
			},
		},
	},
	{
		Title: "关注我们",
		Items: []widget.Link{
			{
				Text: "微信",
				Href: "http://www.ftchinese.com/m/corp/follow.html",
			},
			{
				Text: "微博",
				Href: "http://weibo.com/ftchinese",
			},
			{
				Text: "Linkedin",
				Href: "https://www.linkedin.com/company/4865254?trk=hp-feed-company-name",
			},
			{
				Text: "Facebook",
				Href: "https://www.facebook.com/financialtimeschinese",
			},
			{
				Text: "Twitter",
				Href: "https://twitter.com/FTChinese",
			},
		},
	},
	{
		Title: "FT产品",
		Items: []widget.Link{

			{
				Text: "FT研究院",
				Href: "http://www.ftchinese.com/m/marketing/intelligence.html",
			},
			{
				Text: "FT商学院",
				Href: "http://www.ftchinese.com/channel/mba.html",
			},
			{
				Text: "FT电子书",
				Href: "http://www.ftchinese.com/m/marketing/ebook.html",
			},
			{
				Text: "数据新闻",
				Href: "http://www.ftchinese.com/channel/datanews.html",
			},
			{
				Text: "FT英文版",
				Href: "https://www.ft.com/",
			},
		},
	},
	{
		Title: "移动应用",
		Items: []widget.Link{
			{
				Text: "安卓",
				Href: "http://app.ftchinese.com/androidmobile.html",
			},
		},
	},
}
