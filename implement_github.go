package pkg_login

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

/**
 * Doc : https://docs.github.com/zh/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
 */

const (
	GithubRedirectPath = "https://github.com/login/oauth/authorize"    // Github获取code地址
	GithubTokenPath    = "https://github.com/login/oauth/access_token" // Github获取token地址
	GithubUserInfoPath = "https://api.github.com/user"                 // Github获取用户信息接口
)

func NewGithubConf(id, secret, redirectUrl string) *Config {
	return &Config{
		GithubId:          id,
		GithubSecret:      secret,
		GithubRedirectUrl: redirectUrl,
	}
}

type GithubServer struct {
}

func newGithubServer() *GithubServer {
	return &GithubServer{}
}

func (g *GithubServer) RedirectUrl() (string, error) {
	parsedURL, err := url.Parse(GithubRedirectPath)
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("client_id", config.GithubId)
	queryParams.Add("redirect_uri", config.GithubRedirectUrl)
	queryParams.Add("scope", "user")
	queryParams.Add("state", rand32Str())

	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}

type GithubTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (g *GithubServer) token(code string) (string, error) {
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("client_id", config.GithubId)
	formData.Set("client_secret", config.GithubSecret)
	formData.Set("redirect_uri", config.GithubRedirectUrl)

	headers := map[string]string{"Accept": "application/json", "Content-Type": "application/x-www-form-urlencoded"}
	response, err := postBase(GithubTokenPath, formData.Encode(), headers)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &GithubTokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return "", err
	}

	if len(responseStruct.Error) != 0 {
		return "", errors.New(responseStruct.Error)
	}

	return responseStruct.AccessToken, nil
}

type GithubUserInfo struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Id        int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Message   string `json:"message"`
}

func (g *GithubServer) GetUserinfo(code string) (*Userinfo, error) {
	token, err := g.token(code)
	if err != nil {
		return nil, errors.New("token获取失败:" + err.Error())
	}

	headers := map[string]string{"Authorization": "Bearer " + token}
	response, err := getBase(GithubUserInfoPath, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &GithubUserInfo{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return nil, err
	}

	if len(responseStruct.Message) > 0 {
		return nil, errors.New(responseStruct.Message)
	}

	return &Userinfo{
		Openid:   strconv.Itoa(int(responseStruct.Id)),
		NickName: responseStruct.Name,
		Avatar:   responseStruct.AvatarUrl,
	}, nil
}
