package gonweb

var _ Handler = (*radixTree)(nil)

/*
压缩前缀树
*/
type radixTree struct {
	trees []*radixTrie
}

func MakeradixTree() *radixTree {
	return &radixTree{}
}

func (t *radixTree) Route(method string, pattern string, f GonHandlerFunc) {
	assert1(method != "", "method is empty")
	trie := t.GetTree(method)
	if trie == nil {
		trie = newTrie(method)
		t.trees = append(t.trees, trie)
	}
	trie.addRoute(pattern, f)
}

func (t *radixTree) serve(c *GonContext) GonHandlerFunc {
	trie := t.GetTree(c.Method())
	if trie == nil {
		return nil
	}
	nd := trie.search(c.Path())
	if nd == nil || nd.f == nil {
		return nil
	}
	return nd.f
}

func (t *radixTree) GetTree(method string) *radixTrie {
	for _, tree := range t.trees {
		if tree.method == method {
			return tree
		}
	}
	return nil
}

type radixTrie struct {
	method string
	root   *radixNode
}

func newTrie(m string) *radixTrie {
	return &radixTrie{method: m}
}

func (rt *radixTrie) addRoute(path string, f GonHandlerFunc) {
	if rt.root == nil {
		root := new(radixNode)
		root.fullpath = "/"
		rt.root = root
	}
	rt.root.addRoute(path, f)
}

func (rt *radixTrie) search(p string) *radixNode {
	if rt.root == nil {
		return nil
	}
	root, ok := rt.root.search(p)
	if !ok {
		return nil
	}
	return root
}

type radixNode struct {
	path     string
	fullpath string
	indices  string
	// priority  int
	passCnt int
	// wildChild bool
	nType    int
	end      bool
	children []*radixNode
	f        GonHandlerFunc
}

/*
1.如果当前节点无数据则直接插入
2. 获得公共前缀长度 i 已知 i<=len(p)  i<=len(w)
3. 判断当前p是否需要切分 即 i<len(p)
4. 判断输入w 是否需要切分 即 i<len(w)
5. 都相等则直接修改end进行覆盖操作
*/
func (rn *radixNode) addRoute(word string, f GonHandlerFunc) {
	fullpath := word
	rn.passCnt++

	if rn.path == "" && len(rn.children) == 0 {
		rn.insertChild(word, word, f)
		rn.nType = Root
		return
	}

walk:
	for {
		i := commenPrefixLen(word, rn.path)

		if i < len(rn.path) {
			child := new(radixNode)
			child.path = rn.path[i:]
			child.fullpath = rn.fullpath
			child.end = rn.end
			// child.priority = rn.priority
			// child.wildChild = rn.wildChild
			child.passCnt = rn.passCnt - 1
			child.children = rn.children
			child.nType = Static
			child.f = rn.f
			child.indices = rn.indices

			rn.end = false
			rn.f = nil
			rn.indices = string(rn.path[i])
			rn.fullpath = rn.fullpath[:len(rn.fullpath)-(len(rn.path)-i)]
			rn.path = rn.path[:i]
			rn.children = []*radixNode{child}
		}

		if i < len(word) {
			word = word[i:]
			c := word[0]

			for j := 0; j < len(rn.indices); j++ {
				if rn.indices[j] == c {
					rn = rn.children[j]
					continue walk
				}
			}

			rn.indices += string(c)
			child := &radixNode{
				fullpath: fullpath,
			}
			rn.children = append(rn.children, child)
			rn = child

			rn.insertChild(word, fullpath, f)
		}

		rn.end = true
		rn.f = f
		return
	}

}

/*
1.将当前节点路径p作为前缀，如果不存在该前缀则直接匹配失败
2.在indices中找，找不到则直接失败
3. 重复1，2直到匹配完毕，且最后节点end=true
*/
func (rn *radixNode) search(p string) (n *radixNode, ok bool) {
walk:
	for {
		if rn.path == p {
			if rn.end {
				return rn, true
			}
			return nil, false
		}
		if len(rn.path) > len(p) || p[:len(rn.path)] != rn.path {
			return nil, false
		}
		p = p[len(rn.path):]
		c := p[0]
		for i := 0; i < len(rn.indices); i++ {
			if rn.indices[i] == c {
				rn = rn.children[i]
				continue walk
			}
		}
		return nil, false
	}
}

func commenPrefixLen(w, p string) int {
	i := 0
	for i < len(w) && i < len(p) && w[i] == p[i] {
		i++
	}
	return i
}

func (rn *radixNode) insertChild(p, fp string, f GonHandlerFunc) {
	// for {
	// 	wildcard, i, vlid := findWildcard(p)
	// 	if i < 0 {
	// 		break
	// 	}
	// 	if !vlid {
	// 		panic(fmt.Sprintf("illegal route %s\n", rn.fullpath))
	// 	}

	// 	if len(wildcard) < 2 {
	// 		panic(fmt.Sprintf("illegal route %s\n", rn.fullpath))
	// 	}

	// 	if wildcard[0] == ':' { //param
	// 		if i > 0 {
	// 			rn.path = p[:i]
	// 			p = p[i:]
	// 		}
	// 		child := &radixNode{
	// 			nType:    Param,
	// 			path:     wildcard,
	// 			fullpath: fp,
	// 		}
	// 		rn.addChild(child)
	// 		rn.wildChild = true
	// 		rn = child
	// 		rn.passCnt++

	// 		if len(wildcard) < len(p) {
	// 			p = p[len(wildcard):]
	// 			child := &radixNode{
	// 				priority: 1,
	// 				fullpath: fp,
	// 			}
	// 			rn.addChild(child)
	// 			rn = child
	// 			continue
	// 		}
	// 	}

	// 	i--
	// 	if p[i] != '/' {
	// 		panic(fmt.Sprintf("illegal route %s\n", rn.fullpath))
	// 	}
	// 	rn.path = p[:i]
	// 	//匹配所有空路径
	// 	child := &radixNode{
	// 		wildChild: true,
	// 		nType:     Wildcard,
	// 		fullpath:  fp,
	// 	}
	// 	rn.addChild(child)
	// 	rn.indices = string('/')
	// 	rn = child
	// 	rn.priority++
	// 	child = &radixNode{
	// 		path:     p[i:],
	// 		nType:    Wildcard,
	// 		f:        f,
	// 		fullpath: fp,
	// 	}
	// 	rn.children = []*radixNode{child}
	// 	return

	// }
	rn.path = p
	rn.f = f
	rn.fullpath = fp
}

// 搜索通配符字段，并检查名称中是否有无效字符,wildcard返回的是通配符后面的那一段名字
// *name :name
func findWildcard(p string) (wildcard string, i int, valid bool) {
	for start, c := range []byte(p) {
		if c != ':' && c != '*' {
			continue
		}

		//发现通配符，检查通配符后面的
		for end, c := range []byte(p[start+1:]) {
			switch c {
			case '/':
				return p[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}
		return p[start:], start, valid
	}
	return "", -1, false
}

// 确保通配符子节点处于末尾，因为不会进入indices
func (rn *radixNode) addChild(c *radixNode) {
	// if rn.wildChild && len(rn.children) > 0 {
	// 	wildcardChild := rn.children[len(rn.children)-1]
	// 	rn.children = append(rn.children[:len(rn.children)-1], c, wildcardChild)
	// 	return
	// }
	rn.children = append(rn.children, c)
}
