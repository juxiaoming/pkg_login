package pkg_login

import "fmt"

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Call() {
	fmt.Println("这是登录功能封装")
}
