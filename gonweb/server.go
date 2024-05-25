package gonweb

import (
	"fmt"
	"net"
	"net/http"
)

type GonHandlerFunc func(ctx *GonContext)

type Routeable interface {
	Route(method, pattern string, f GonHandlerFunc)
}

type Server interface {
	http.Handler
	Routeable
	Start() error
	ShutDown() error
}

type WebServer struct {
	name    string
	addr    string
	handler Handler
}

var _ Server = (*WebServer)(nil)

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := MakeContext(w, r)
	f := s.handler.serve(c)
	if f == nil {
		c.PageNotFound()
		return
	}
	f(c)
}

func (s *WebServer) Route(method, pattern string, f GonHandlerFunc) {
	s.handler.Route(method, pattern, f)
}

func (s *WebServer) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(fmt.Sprintf("server-%s start failed:%s\n", s.name, err))
	}
	fmt.Printf("server-%s start at port:%s", s.name, s.addr)
	//TODO服务注册与发现 map[xxx]xxx append
	return http.Serve(l, s)
}

func (s *WebServer) ShutDown() error {
	fmt.Printf("server-%s is Closing", s.name)
	return nil
}

func MakeWebServer(name, addr string, h Handler) *WebServer {
	if h == nil {
		h = &RouteBaseOnTree{
			trees: make(map[string]*treeNode),
		}
	}
	return &WebServer{
		name:    name,
		addr:    addr,
		handler: h,
	}
}
