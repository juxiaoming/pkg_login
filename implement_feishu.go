package pkg_login

import (
	"encoding/json"
	"errors"
	"net/url"
)

/**
 * Doc : https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/login-overview
 */

const (
	FeiShuRedirectPath = "https://passport.feishu.cn/suite/passport/oauth/authorize" // 飞书获取code地址
	FeiShuTokenPath    = "https://passport.feishu.cn/suite/passport/oauth/token"     // 飞书获取token地址
	FeiShuUserInfoPath = "https://passport.feishu.cn/suite/passport/oauth/userinfo"  // 飞书获取用户信息接口
)

func NewFeiShuConf(id, secret, redirectUrl string) *Config {
	return &Config{
		FeiShuId:          id,
		FeiShuSecret:      secret,
		FeiShuRedirectUrl: redirectUrl,
	}
}

type FeiShuServer struct {
}

func newFeiShuServer() *FeiShuServer {
	return &FeiShuServer{}
}

func (f *FeiShuServer) RedirectUrl() (string, error) {
	parsedURL, err := url.Parse(FeiShuRedirectPath)
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("redirect_uri", config.FeiShuRedirectUrl)
	queryParams.Add("client_id", config.FeiShuId)
	queryParams.Add("response_type", "code")
	queryParams.Add("state", rand32Str())

	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}

type FeiShuTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (f *FeiShuServer) token(code string) (string, error) {
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("client_id", config.FeiShuId)
	formData.Set("client_secret", config.FeiShuSecret)
	formData.Set("redirect_uri", config.FeiShuRedirectUrl)
	formData.Set("grant_type", "authorization_code")

	headers := map[string]string{"Accept": "application/json", "Content-Type": "application/x-www-form-urlencoded"}
	response, err := postBase(FeiShuTokenPath, formData.Encode(), headers)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &FeiShuTokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return "", err
	}

	if len(responseStruct.Error) != 0 {
		return "", errors.New(responseStruct.ErrorDescription)
	}

	return responseStruct.AccessToken, nil
}

type FeiShuUserInfo struct {
	Sub          string `json:"sub"`
	Picture      string `json:"picture"`
	Name         string `json:"name"`
	EnName       string `json:"en_name"`
	TenantKey    string `json:"tenant_key"`
	AvatarUrl    string `json:"avatar_url"`
	AvatarThumb  string `json:"avatar_thumb"`
	AvatarMiddle string `json:"avatar_middle"`
	AvatarBig    string `json:"avatar_big"`
	OpenId       string `json:"open_id"`
	UnionId      string `json:"union_id"`
	UserId       string `json:"user_id"`
	Mobile       string `json:"mobile"`
	Message      string `json:"message"`
}

func (f *FeiShuServer) GetUserinfo(code string) (*Userinfo, error) {
	token, err := f.token(code)
	if err != nil {
		return nil, errors.New("token获取失败:" + err.Error())
	}

	headers := map[string]string{"Authorization": "Bearer " + token}
	response, err := getBase(FeiShuUserInfoPath, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &FeiShuUserInfo{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return nil, err
	}

	if len(responseStruct.Message) > 0 {
		return nil, errors.New(responseStruct.Message)
	}

	return &Userinfo{
		Openid:   responseStruct.OpenId,
		UnionId:  responseStruct.UnionId,
		NickName: responseStruct.Name,
		Avatar:   responseStruct.AvatarUrl,
		Mobile:   responseStruct.Mobile,
	}, nil
}
