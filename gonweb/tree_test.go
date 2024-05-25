package gonweb

import (
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

func TestSort(t *testing.T) {
	a := []int{1}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	t.Log(a)
}
