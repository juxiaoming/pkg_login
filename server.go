package pkg_login

import "errors"

const (
	ImplementGoogle   int8 = 1 // 谷歌
	ImplementWeiXin   int8 = 2 // 微信
	ImplementGithub   int8 = 3 // github
	ImplementQq       int8 = 4 // qq
	ImplementWeiBo    int8 = 5 // 微博
	ImplementDingDing int8 = 6 // 钉钉
	ImplementGitee    int8 = 7 // gitee码云
	ImplementFeiShu   int8 = 8 // 飞书
)

var (
	config  Config // 全局配置
	hasInit bool   // 配置是否初始化
)

type Config struct {
	GoogleId            string `json:"google_id"`
	GoogleSecret        string `json:"google_secret"`
	GoogleRedirectUrl   string `json:"google_redirect_url"`
	GithubId            string `json:"github_id"`
	GithubSecret        string `json:"github_secret"`
	GithubRedirectUrl   string `json:"github_redirect_url"`
	GiteeId             string `json:"gitee_id"`
	GiteeSecret         string `json:"gitee_secret"`
	GiteeRedirectUrl    string `json:"gitee_redirect_url"`
	DingDingId          string `json:"ding_ding_id"`
	DingDingSecret      string `json:"ding_ding_secret"`
	DingDingRedirectUrl string `json:"ding_ding_redirect_url"`
	FeiShuId            string `json:"fei_shu_id"`
	FeiShuSecret        string `json:"fei_shu_secret"`
	FeiShuRedirectUrl   string `json:"fei_shu_redirect_url"`
}

type Userinfo struct {
	Openid   string `json:"openid"`
	UnionId  string `json:"unionId"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
}

type Ability interface {
	RedirectUrl() (string, error)
	GetUserinfo(code string) (*Userinfo, error)
}

type Server struct {
	client      Ability
	ImplementId int8 `json:"implement_id"`
}

func Init(conf *Config) {
	config = *conf
	hasInit = true
}

func NewServer(implementId int8) (*Server, error) {
	if !hasInit {
		return nil, errors.New("配置未初始化,请先调用【Init】方法")
	}

	var client Ability
	switch implementId {
	case ImplementGoogle:
		if len(config.GoogleId) == 0 || len(config.GoogleSecret) == 0 || len(config.GoogleRedirectUrl) == 0 {
			return nil, errors.New("缺失配置文件")
		}
		client = newGoogleServer()
	case ImplementGithub:
		if len(config.GithubId) == 0 || len(config.GithubSecret) == 0 || len(config.GithubRedirectUrl) == 0 {
			return nil, errors.New("缺失配置文件")
		}
		client = newGithubServer()
	case ImplementGitee:
		if len(config.GiteeId) == 0 || len(config.GiteeSecret) == 0 || len(config.GiteeRedirectUrl) == 0 {
			return nil, errors.New("缺失配置文件")
		}
		client = newGiteeServer()
	case ImplementDingDing:
		if len(config.DingDingId) == 0 || len(config.DingDingSecret) == 0 || len(config.DingDingRedirectUrl) == 0 {
			return nil, errors.New("缺失配置文件")
		}
		client = newDingDingServer()
	case ImplementFeiShu:
		if len(config.FeiShuId) == 0 || len(config.FeiShuSecret) == 0 || len(config.FeiShuRedirectUrl) == 0 {
			return nil, errors.New("缺失配置文件")
		}
		client = newFeiShuServer()
	default:
		return nil, errors.New("未定义实现")
	}

	return &Server{
		client:      client,
		ImplementId: implementId,
	}, nil
}

func (s *Server) RedirectUrl() (string, error) {
	return s.client.RedirectUrl()
}

func (s *Server) GetUserinfo(code string) (*Userinfo, error) {
	return s.client.GetUserinfo(code)
}
