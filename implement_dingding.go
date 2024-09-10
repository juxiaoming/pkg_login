package pkg_login

import (
	"encoding/json"
	"errors"
	"net/url"
)

/**
 * Doc : https://open.dingtalk.com/document/orgapp/tutorial-obtaining-user-personal-information?spm=ding_open_doc.document.0.0.6a584a97pcBXIV
 */

const (
	DingDingRedirectPath = "https://login.dingtalk.com/oauth2/auth"               // 钉钉获取code地址
	DingDingTokenPath    = "https://api.dingtalk.com/v1.0/oauth2/userAccessToken" // 钉钉获取token地址
	DingDingUserInfoPath = "https://api.dingtalk.com/v1.0/contact/users/me"       // 钉钉获取用户信息接口
)

func NewDingDingConf(id, secret, redirectUrl string) *Config {
	return &Config{
		DingDingId:          id,
		DingDingSecret:      secret,
		DingDingRedirectUrl: redirectUrl,
	}
}

type DingDingServer struct {
}

func newDingDingServer() *DingDingServer {
	return &DingDingServer{}
}

func (d *DingDingServer) RedirectUrl() (string, error) {
	parsedURL, err := url.Parse(DingDingRedirectPath)
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("redirect_uri", config.DingDingRedirectUrl)
	queryParams.Add("client_id", config.DingDingId)
	queryParams.Add("response_type", "code")
	queryParams.Add("scope", "openid")
	queryParams.Add("state", "authCode")
	queryParams.Add("prompt", "consent")

	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}

type DingDingTokenResponse struct {
	ExpireIn     int    `json:"expireIn"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Message      string `json:"message"`
}

func (d *DingDingServer) token(code string) (string, error) {
	payload := map[string]string{
		"clientId":     config.DingDingId,
		"clientSecret": config.DingDingSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}

	payloadBytes, _ := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json"}
	response, err := postBase(DingDingTokenPath, string(payloadBytes), headers)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &DingDingTokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return "", err
	}

	if len(responseStruct.Message) != 0 {
		return "", errors.New(responseStruct.Message)
	}

	return responseStruct.AccessToken, nil
}

type DingDingUserInfo struct {
	Nick      string `json:"nick"`
	UnionId   string `json:"unionId"`
	AvatarUrl string `json:"avatarUrl"`
	OpenId    string `json:"openId"`
	Mobile    string `json:"mobile"`
	StateCode string `json:"stateCode"`
	Visitor   bool   `json:"visitor"`
	Message   string `json:"message"`
}

func (d *DingDingServer) GetUserinfo(code string) (*Userinfo, error) {
	token, err := d.token(code)
	if err != nil {
		return nil, errors.New("token获取失败:" + err.Error())
	}

	headers := map[string]string{"x-acs-dingtalk-access-token": token}
	response, err := getBase(DingDingUserInfoPath, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &DingDingUserInfo{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return nil, err
	}

	if len(responseStruct.Message) > 0 {
		return nil, errors.New(responseStruct.Message)
	}

	return &Userinfo{
		Openid:   responseStruct.OpenId,
		UnionId:  responseStruct.UnionId,
		NickName: responseStruct.Nick,
		Avatar:   responseStruct.AvatarUrl,
		Mobile:   responseStruct.Mobile,
	}, nil
}
