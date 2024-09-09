package pkg_login

import (
	"encoding/json"
	"errors"
	"fmt"
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
	payload := map[string]string{
		"code":          code,
		"client_id":     config.GithubId,
		"client_secret": config.GithubSecret,
		"redirect_uri":  config.GithubRedirectUrl,
	}
	payloadBytes, _ := json.Marshal(payload)

	headers := map[string]string{"Accept": "application/json"}

	response, err := postBase(GithubTokenPath, string(payloadBytes), headers)
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

	if len(responseStruct.ErrorDescription) != 0 {
		return "", errors.New(responseStruct.ErrorDescription)
	}

	fmt.Println("token响应数据", responseStruct)

	return responseStruct.AccessToken, nil
}

type GithubUserInfo struct {
	Login     string `json:"login"`
	Id        int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Message   string `json:"message"`
}

func (g *GithubServer) GetUserinfo(code string) (*Userinfo, error) {
	token, err := g.token(code)
	if err != nil {
		return nil, err
	}

	fmt.Println("这是token数据", token)

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
		NickName: responseStruct.Login,
		Avatar:   responseStruct.AvatarUrl,
	}, nil
}
