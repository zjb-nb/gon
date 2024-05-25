package middleware

import "gon/gonweb"

type ResponsePageMiddleware struct {
	resps map[int][]byte
}

func NewRespMiddleware() *ResponsePageMiddleware {
	return &ResponsePageMiddleware{
		resps: map[int][]byte{},
	}
}

func (m *ResponsePageMiddleware) SendRespBuilder(next gonweb.GonHandlerFunc) gonweb.GonHandlerFunc {
	return func(ctx *gonweb.GonContext) {
		next(ctx)
		ctx.W.WriteHeader(ctx.ResponStatus)
		ctx.W.Write(ctx.ResponsData)
	}
}

func (m *ResponsePageMiddleware) HookBuilder(next gonweb.GonHandlerFunc) gonweb.GonHandlerFunc {
	return func(ctx *gonweb.GonContext) {
		next(ctx)
		//todo 配合template
		resp, ok := m.resps[ctx.ResponStatus]
		if ok {
			ctx.ResponsData = resp
		}
	}
}

func (m *ResponsePageMiddleware) Addpage(code int, page []byte) {
	m.resps[code] = page
}
