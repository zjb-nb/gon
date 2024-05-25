package gonweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type GonContext struct {
	W            http.ResponseWriter
	R            *http.Request
	Values       map[string]string
	ResponsData  []byte
	ResponStatus int
}

type ResponsData struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (g *GonContext) Write(b []byte) (int, error) {
	return g.W.Write(b)
}

func (g *GonContext) Method() string {
	return g.R.Method
}
func (g *GonContext) Path() string {
	return g.R.URL.Path
}

func (g *GonContext) WriteHeader(code int) {
	g.W.WriteHeader(code)
}

func (g *GonContext) ResponsJson(code int, data string) {
	d := &ResponsData{
		Code: code,
		Data: data,
	}
	json_data, err := json.Marshal(d)
	if err != nil {
		g.ResponStatus = http.StatusInternalServerError
		g.ResponsData = []byte(err.Error())
		return
	}
	g.WriteHeader(code)
	g.Write(json_data)
}

func (g *GonContext) HeaderSet(k, v string) {
	g.W.Header().Set(k, v)
}

func (g *GonContext) Respons(code int, data string) {
	g.ResponStatus = code
	g.ResponsData = []byte(data)
}

func (g *GonContext) JsonOk(data string) {
	g.ResponsJson(http.StatusOK, data)
}

func (g *GonContext) Ok(data string) {
	g.Respons(http.StatusOK, data)
}
func (g *GonContext) JsonPageNotFound() {
	g.ResponsJson(http.StatusNotFound, "page not found")
}
func (g *GonContext) PageNotFound() {
	g.Respons(http.StatusNotFound, "page not found")
}
func (g *GonContext) JsonServerError(data string) {
	g.ResponsJson(http.StatusInternalServerError, data)
}
func (g *GonContext) ServerError(data string) {
	g.Respons(http.StatusInternalServerError, data)
}

func (g *GonContext) Form() url.Values {

	err := g.R.ParseForm()
	if err != nil {
		return nil
	}
	return g.R.PostForm

}

func (g *GonContext) FormValue(key string) (string, error) {
	if err := g.R.ParseForm(); err != nil {
		return "", err
	}
	return g.R.FormValue(key), nil
}

func (g *GonContext) QueryParam() url.Values {
	return g.R.URL.Query()
}

func (g *GonContext) QueryValue(key string) (string, error) {
	p := g.QueryParam()
	if p == nil {
		return "", errors.New("not found res")
	}
	vals, ok := p[key]
	if !ok {
		return "", errors.New("not found this key")
	}
	return vals[0], nil
}

func MakeContext(w http.ResponseWriter, r *http.Request) *GonContext {
	return &GonContext{
		W:            w,
		R:            r,
		Values:       make(map[string]string),
		ResponStatus: http.StatusOK,
	}
}
