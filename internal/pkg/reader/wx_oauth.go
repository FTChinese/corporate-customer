package reader

type WxOAuthCodeParams struct {
	AppID        string `json:"appId"`
	RedirectURI  string `json:"redirectUri"`
	ResponseType string `json:"responseType"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}
