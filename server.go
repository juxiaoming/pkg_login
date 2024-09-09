package pkg_login

import "errors"

const (
	ImplementGoogle   int8 = 1
	ImplementWeiXin   int8 = 2
	ImplementGithub   int8 = 3
	ImplementQq       int8 = 4
	ImplementWeiBo    int8 = 5
	ImplementDingDing int8 = 6
)

var (
	config  Config // 全局配置
	hasInit bool   // 配置是否初始化
)

type Config struct {
	GoogleId          string `json:"google_id"`
	GoogleSecret      string `json:"google_secret"`
	GoogleRedirectUrl string `json:"google_redirect_url"`
}

type Userinfo struct {
	Openid   string `json:"openid"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

type Ability interface {
	RedirectUrl() (string, error)
	GetUserinfo(code string) (*Userinfo, error)
}

type Server struct {
	client      Ability
	ImplementId int8 `json:"implement_id"`
}

func Init(conf Config) {
	config = conf
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
