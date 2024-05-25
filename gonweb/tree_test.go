package gonweb

import (
	"log"
	"net/http"
	"sort"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestTreeRoute(t *testing.T) {
	// 非法路由 /a/**  /a*b
	require.Panics(t, func() {
		t := NewTree()
		t.Route("GET", "/a/**", nil)
	})
	require.Panics(t, func() {
		tree := NewTree()
		tree.Route("GET", "/a*b", nil)
	})

	tree := NewTree()
	Tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "get-/",
			method: "GET",
			path:   "/",
		},
		{
			name:   "get-/home",
			method: "GET",
			path:   "/home",
		},
		{
			name:   "get-/home/user",
			method: "GET",
			path:   "/home/user",
		},
		{
			name:   "get-/home/index",
			method: "GET",
			path:   "/home/index",
		},
		{
			name:   "post-/home/login",
			method: "POST",
			path:   "/home/login",
		},
	}

	for _, tt := range Tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tree.route(tt.method, tt.path)
			assert.Equal(t, res.fullpath, tt.path, tt.name)
		})
	}
}

func TestMain(t *testing.T) {
	s := MakeWebServer("web", ":8080", nil)
	s.Route("GET", "/", func(ctx *GonContext) { ctx.JsonOk("hi") })
	s.Route("GET", "/home", func(ctx *GonContext) { ctx.Ok("home") })
	s.Route("GET", "/home/user", func(ctx *GonContext) { ctx.Ok("home/user") })
	s.Route("GET", "/login", func(ctx *GonContext) { ctx.Ok("login") })
	s.Route("GET", "/*", func(ctx *GonContext) { ctx.Ok("/* wildcard") })
	s.Route("GET", "/home/*", func(ctx *GonContext) { ctx.Ok("/home/* wildcard") })
	s.Route("GET", "/home/user/:id", func(ctx *GonContext) { ctx.Ok("/home/user/:id ") })
	s.Route("GET", "/home/~^ab$", func(ctx *GonContext) { ctx.Ok("/home/~^ab$") })

	s.Route("GET", "/home/user/index", func(ctx *GonContext) { ctx.JsonOk("post-hi") })

	go func() {
		log.Fatal(s.Start())
	}()
	res, err := http.Get("http://localhost:8080/")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get("http://localhost:8080/sss")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get("http://localhost:8080/home/login/index")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res, err = http.Get("http://localhost:8080/home/user/index")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	res, err = http.Get("http://localhost:8080/home/user/ss")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

}

func TestSort(t *testing.T) {
	a := []int{1}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	t.Log(a)
}
