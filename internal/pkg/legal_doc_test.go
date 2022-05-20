package pkg

import (
	"testing"
)

func TestLegalDoc_Rendered(t *testing.T) {

	md := `
###### Cookies的使用<a name="cookie"></a>

我们利用 cookies追访您的喜好，让您更个性化地使用我们的网站/APP。Cookie 即是网站/APP发送到您电脑硬盘上的一条信息，以便网站/APP记忆您的身份。该信息可能包涵您使用网站/APP的相关信息、您的电脑IP地址、PC/移动设备类型、基本资料、消费习惯等等。此外，如果您是通过第三方网站链接到我方网站，则该信息中还包括链接页面的URL。如果您是注册用户，您所注册的详细资料也将包括在内，以便查证。

我们可能出于以下目的利用从 cookies上得到的信息：

* 辨别返回用户、登录者和注册用户，并为注册用户呈现个性化的页面。
* 使您更自如地浏览我们的网站和使用我们的APP
* 订阅付费账户使用参考
* 追踪您对网站APP的访问和使用情况，以便根据您的需求进一步完善网站/APP和向您提供更好的服务
* 核查无效或违法违规和违反用户协议的账户使用；就产品服务等更新与您沟通
* 向您提供更针对您需求的定制化广告推送
* 帮助处理申请流程，如您通过我们的网站/APP申请和应聘`

	type fields struct {
		LegalTeaser LegalTeaser
		Body        string
	}
	tests := []struct {
		name   string
		fields fields
		want   LegalDoc
	}{
		{
			name: "Render",
			fields: fields{
				LegalTeaser: LegalTeaser{},
				Body:        md,
			},
			want: LegalDoc{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LegalDoc{
				LegalTeaser: tt.fields.LegalTeaser,
				Body:        tt.fields.Body,
			}

			got := l.Rendered()

			//if got := l.Rendered(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Rendered() = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", got.Body)
		})
	}
}
