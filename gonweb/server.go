package gonweb

import (
	"fmt"
	"net"
	"net/http"
	"time"
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
	name     string
	addr     string
	handler  Handler
	builders []FilterBuilder
}

var _ Server = (*WebServer)(nil)

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := MakeContext(w, r)
	f := s.handler.serve(c)
	if f == nil {
		f = func(ctx *GonContext) { ctx.PageNotFound() }
	}

	for i := len(s.builders) - 1; i >= 0; i-- {
		f = s.builders[i](f)
	}
	f(c)
}

func (s *WebServer) Route(method, pattern string, f GonHandlerFunc) {
	s.handler.Route(method, pattern, f)
}

func (s *WebServer) GET(pattern string, f GonHandlerFunc) {
	s.Route("GET", pattern, f)
}

func (s *WebServer) POST(pattern string, f GonHandlerFunc) {
	s.Route("POST", pattern, f)
}
func (s *WebServer) PUT(pattern string, f GonHandlerFunc) {
	s.Route("PUT", pattern, f)
}
func (s *WebServer) DELETE(pattern string, f GonHandlerFunc) {
	s.Route("DELETE", pattern, f)
}
func (s *WebServer) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(fmt.Sprintf("server-%s start failed:%s\n", s.name, err))
	}
	fmt.Printf("server-%s start at port:%s\n", s.name, s.addr)
	//TODO服务注册与发现 map[xxx]xxx append
	return http.Serve(l, s)
}

func (s *WebServer) ShutDown() error {
	time.Sleep(time.Second * 5)
	fmt.Printf("server-%s start Closing\n", s.name)
	return nil
}

func MakeWebServer(name, addr string, h Handler, builders ...FilterBuilder) *WebServer {
	if h == nil {
		h = &RouteBaseOnTree{
			trees: make(map[string]*treeNode),
		}
	}
	return &WebServer{
		name:     name,
		addr:     addr,
		handler:  h,
		builders: builders,
	}
}

type RouteGroup struct {
	prefix string
	s      Server
	h      []FilterBuilder
}

func NewGroup(p string, s Server, b ...FilterBuilder) *RouteGroup {
	return &RouteGroup{
		prefix: p,
		s:      s,
		h:      b,
	}
}

func (g *RouteGroup) Route(method, pattern string, f GonHandlerFunc) {
	g.s.Route(method, g.prefix+pattern, f)
}
