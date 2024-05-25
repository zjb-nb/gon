package gonweb

type Handler interface {
	Routeable
	serve(*GonContext) GonHandlerFunc
}

var _ Handler = (*RouteBaseOnMap)(nil)

type RouteBaseOnMap struct {
	m map[string]GonHandlerFunc
}

func (r *RouteBaseOnMap) Route(method string, path string, f GonHandlerFunc) {
	r.m[method+"#"+path] = f
}
func (r *RouteBaseOnMap) serve(c *GonContext) GonHandlerFunc {
	f, ok := r.m[c.Method()+"#"+c.Path()]
	if !ok {
		return nil
	}
	return f
}

func MakeMapRoute() *RouteBaseOnMap {
	return &RouteBaseOnMap{
		m: make(map[string]GonHandlerFunc),
	}
}
