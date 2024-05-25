package gonweb

import (
	"encoding/json"
	"net/http"
)

type GonContext struct {
	W      http.ResponseWriter
	R      *http.Request
	Values map[string]string
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
		g.WriteHeader(http.StatusInternalServerError)
		g.Write([]byte(err.Error()))
		return
	}
	g.WriteHeader(code)
	g.Write(json_data)
}

func (g *GonContext) HeaderSet(k, v string) {
	g.W.Header().Set(k, v)
}

func (g *GonContext) Respons(code int, data string) {
	g.WriteHeader(code)
	g.Write([]byte(data))
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
func MakeContext(w http.ResponseWriter, r *http.Request) *GonContext {
	return &GonContext{
		W:      w,
		R:      r,
		Values: make(map[string]string),
	}
}
