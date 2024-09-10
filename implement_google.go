package pkg_login

import (
	"encoding/json"
	"errors"
	"net/url"
)

/**
 * Doc : https://developers.google.com/identity/protocols/oauth2/web-server?hl=zh-cn
 */

const (
	GoogleRedirectPath = "https://accounts.google.com/o/oauth2/auth"     // 谷歌获取code地址
	GoogleTokenPath    = "https://oauth2.googleapis.com/token"           // 谷歌获取token地址
	GoogleUserInfoPath = "https://www.googleapis.com/oauth2/v2/userinfo" // 谷歌获取用户信息接口
)

func NewGoogleConf(id, secret, redirectUrl string) *Config {
	return &Config{
		GoogleId:          id,
		GoogleSecret:      secret,
		GoogleRedirectUrl: redirectUrl,
	}
}

type GoogleServer struct {
}

func newGoogleServer() *GoogleServer {
	return &GoogleServer{}
}

func (g *GoogleServer) RedirectUrl() (string, error) {
	parsedURL, err := url.Parse(GoogleRedirectPath)
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("response_type", "code")
	queryParams.Add("client_id", config.GoogleId)
	queryParams.Add("redirect_uri", config.GoogleRedirectUrl)
	queryParams.Add("scope", "https://www.googleapis.com/auth/userinfo.profile")
	queryParams.Add("access_type", "offline")

	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}

type GoogleTokenResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
}

func (g *GoogleServer) token(code string) (string, error) {
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("client_id", config.GoogleId)
	formData.Set("client_secret", config.GoogleSecret)
	formData.Set("redirect_uri", config.GoogleRedirectUrl)
	formData.Set("grant_type", "authorization_code")
	headers := map[string]string{"Accept": "application/json", "Content-Type": "application/x-www-form-urlencoded"}
	response, err := postBase(GoogleTokenPath, formData.Encode(), headers)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &GoogleTokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return "", err
	}

	if len(responseStruct.ErrorDescription) != 0 {
		return "", errors.New(responseStruct.ErrorDescription)
	}

	return responseStruct.AccessToken, nil
}

type GoogleUserInfo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Error   struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (g *GoogleServer) GetUserinfo(code string) (*Userinfo, error) {
	token, err := g.token(code)
	if err != nil {
		return nil, errors.New("token获取失败:" + err.Error())
	}

	headers := map[string]string{"Authorization": "Bearer " + token}
	response, err := getBase(GoogleUserInfoPath, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &GoogleUserInfo{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return nil, err
	}

	if responseStruct.Error.Code != 0 {
		return nil, errors.New(responseStruct.Error.Message)
	}

	return &Userinfo{
		Openid:   responseStruct.Id,
		NickName: responseStruct.Name,
		Avatar:   responseStruct.Picture,
	}, nil
}
