package gonweb

import (
	"regexp"
	"sort"
	"strings"
)

const (
	Root = iota
	Wildcard
	Param
	Regular
	Static
)

var _ (Handler) = (*RouteBaseOnTree)(nil)

type RouteBaseOnTree struct {
	trees map[string]*treeNode
}

/*
1.验证合法性
2. 根据方法查找对应树，作为根节点
3. 根据子路径查找或添加子节点，作为根节点
4. 重复3直到子路径消耗完毕
*/
func (t *RouteBaseOnTree) Route(method, pattern string, f GonHandlerFunc) {
	root := t.route(method, pattern)
	root.f = f
}

func (t *RouteBaseOnTree) route(method, pattern string) *treeNode {
	root, ok := t.trees[method]
	if !ok {
		root = NewRootNode()
		t.trees[method] = root
	}
	pattern = strings.Trim(strings.Trim(pattern, ""), "/")
	ok = t.validPath(pattern)
	if !ok {
		panic("add illegal route!")
	}
	if pattern == "" {
		return root
	}
	paths := strings.Split(pattern, "/")
	fullpath := ""
	for _, path := range paths {
		fullpath = fullpath + "/" + path
		root = t.findOrCreate(path, fullpath, root)
	}
	return root
}

/*
通配符必须作为最后出现且前一个必须为/
*/
func (t *RouteBaseOnTree) validPath(p string) bool {
	pos := strings.Index(p, "*")
	if pos > 0 {
		if p[pos-1] != '/' || pos != len(p)-1 {
			return false
		}
	}
	return true
}

func (t *RouteBaseOnTree) findOrCreate(p string, fullpath string, root *treeNode) *treeNode {
	for _, n := range root.children {
		if n.path == p {
			return n
		}
	}
	var newnode *treeNode
	if p == "*" {
		newnode = NewWildCardNode(fullpath)
	} else if p[0] == ':' {
		newnode = NewParamNode(p, fullpath)
	} else if p[0] == '~' {
		newnode = NewRegularNode(p, fullpath)
	} else {
		newnode = NewStaticNode(p, fullpath)
	}
	root.children = append(root.children, newnode)
	return newnode
}

func (t *RouteBaseOnTree) serve(c *GonContext) GonHandlerFunc {
	node := t.find(c)
	if node == nil {
		return nil
	}
	return node.f
}

func (t *RouteBaseOnTree) find(c *GonContext) *treeNode {
	root, ok := t.trees[c.Method()]
	if !ok {
		return nil
	}
	pattern := strings.Trim(strings.Trim(c.Path(), ""), "/")
	if pattern == "" {
		return root
	}
	paths := strings.Split(pattern, "/")
	for _, path := range paths {
		root, ok = t.findNode(path, root, c)
		if !ok {
			return nil
		}
	}
	return root
}

func (t *RouteBaseOnTree) findNode(p string, n *treeNode, c *GonContext) (*treeNode, bool) {
	var res = make([]*treeNode, 0, 2)
	for _, child := range n.children {
		if child.match(p) {
			res = append(res, child)
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].nodetype < res[j].nodetype
	})
	node := res[len(res)-1]
	if node.nodetype == Param {
		c.Values[node.path[1:]] = p
	}
	return node, true
}

type treeNode struct {
	path     string
	fullpath string
	match    func(p string) bool
	nodetype int
	f        GonHandlerFunc
	children []*treeNode
}

func NewRootNode() *treeNode {
	return &treeNode{
		path:     "/",
		fullpath: "/",
		nodetype: Root,
		match:    func(p string) bool { return true },
	}
}

func NewWildCardNode(fullpath string) *treeNode {
	return &treeNode{
		path:     "*",
		fullpath: fullpath,
		nodetype: Wildcard,
		match: func(p string) bool {
			return true
		},
	}
}

func NewStaticNode(path, fullpath string) *treeNode {
	return &treeNode{
		path:     path,
		fullpath: fullpath,
		nodetype: Static,
		match:    func(p string) bool { return p == path },
	}
}

func NewParamNode(path, fullpath string) *treeNode {
	return &treeNode{
		path:     path,
		fullpath: fullpath,
		nodetype: Param,
		match: func(p string) bool {
			return p != "*"
		},
	}
}

func NewRegularNode(path, fullpath string) *treeNode {
	return &treeNode{
		path:     path,
		fullpath: fullpath,
		nodetype: Static,
		match: func(p string) bool {
			ruler := regexp.MustCompile(path[1:])
			return ruler.Match([]byte(p))
		},
	}
}

func NewTree() *RouteBaseOnTree {
	return &RouteBaseOnTree{
		trees: make(map[string]*treeNode),
	}
}
