package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"net/http"
)

// WxLogin performs the Step 2 and Step 3 of OAuth workflow as described by
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317851&token=&lang=zh_CN.
// or
// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
//
// Request body:
// * code: xxxxx
//
// Request header:
// * X-App-Id: xxxx
func (c Client) WxLogin(appID string, body io.Reader, client http.Header) (*http.Response, error) {
	url := c.baseURL + pathWxLogin

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(client).
		SetHeader(keyWxAppID, appID).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// WxRefresh allows user to refresh userinfo.
// Request header
//* `X-App-Id`
// and client metadata.
//
// Request body
// * sessionId: string
func (c Client) WxRefresh(appID string, body io.Reader, client http.Header) (*http.Response, error) {
	url := c.baseURL + pathWxRefresh

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(client).
		SetHeader(keyWxAppID, appID).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
