package gonweb

import "strings"

/*
通配符必须作为最后出现且前一个必须为/
*/
func validPath(p string) bool {
	pos := strings.Index(p, "*")
	if pos > 0 {
		if p[pos-1] != '/' || pos != len(p)-1 {
			return false
		}
	}
	return true
}

func assert1(condition bool, text string) {
	if !condition {
		panic(text)
	}
}
