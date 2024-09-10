package pkg_login

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

/**
 * Doc : https://gitee.com/api/v5/oauth_doc#/
 */

const (
	GiteeRedirectPath = "https://gitee.com/oauth/authorize" // Gitee获取code地址
	GiteeTokenPath    = "https://gitee.com/oauth/token"     // Gitee获取token地址
	GiteeUserInfoPath = "https://gitee.com/api/v5/user"     // Gitee获取用户信息接口
)

func NewGiteeConf(id, secret, redirectUrl string) *Config {
	return &Config{
		GiteeId:          id,
		GiteeSecret:      secret,
		GiteeRedirectUrl: redirectUrl,
	}
}

type GiteeServer struct {
}

func newGiteeServer() *GiteeServer {
	return &GiteeServer{}
}

func (g *GiteeServer) RedirectUrl() (string, error) {
	parsedURL, err := url.Parse(GiteeRedirectPath)
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("client_id", config.GiteeId)
	queryParams.Add("redirect_uri", config.GiteeRedirectUrl)
	queryParams.Add("response_type", "code")

	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}

type GiteeTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	CreatedAt        int    `json:"created_at"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (g *GiteeServer) token(code string) (string, error) {
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("client_id", config.GiteeId)
	formData.Set("client_secret", config.GiteeSecret)
	formData.Set("redirect_uri", config.GiteeRedirectUrl)
	formData.Set("grant_type", "authorization_code")

	headers := map[string]string{"Accept": "application/json", "Content-Type": "application/x-www-form-urlencoded"}
	response, err := postBase(GiteeTokenPath, formData.Encode(), headers)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &GiteeTokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return "", err
	}

	if len(responseStruct.Error) != 0 {
		return "", errors.New(responseStruct.ErrorDescription)
	}

	return responseStruct.AccessToken, nil
}

type GiteeUserInfo struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Id        int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Message   string `json:"message"`
}

func (g *GiteeServer) GetUserinfo(code string) (*Userinfo, error) {
	token, err := g.token(code)
	if err != nil {
		return nil, errors.New("token获取失败:" + err.Error())
	}

	parsedURL, err := url.Parse(GiteeUserInfoPath)
	if err != nil {
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Add("access_token", token)
	parsedURL.RawQuery = queryParams.Encode()

	response, err := getBase(parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseStruct := &GiteeUserInfo{}
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
