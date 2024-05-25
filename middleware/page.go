package middleware

import (
	"bytes"
	"html/template"
)

func NotFoundPage() []byte {
	page := `<html><h1>404 NOT FOUND</h1></html>`
	tpl, err := template.New("404").Parse(page)
	if err != nil {
		panic("parser notfoundpage failed!")
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, nil)
	if err != nil {
		panic("template exec failed!")
	}
	return buf.Bytes()
}
