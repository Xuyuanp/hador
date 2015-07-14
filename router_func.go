package hador

import (
	"github.com/Xuyuanp/hador/swagger"
)

type RouterFunc func(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf

func (r RouterFunc) AddRoute(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r(method, pattern, handler, filters...)
}

func (r RouterFunc) Get(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(GET, pattern, handler, filters...)
}

func (r RouterFunc) Post(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(POST, pattern, handler, filters...)
}

func (r RouterFunc) Put(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(PUT, pattern, handler, filters...)
}

func (r RouterFunc) Delete(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(DELETE, pattern, handler, filters...)
}

func (r RouterFunc) Patch(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(PATCH, pattern, handler, filters...)
}

func (r RouterFunc) Trace(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(TRACE, pattern, handler, filters...)
}

func (r RouterFunc) Connect(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(CONNECT, pattern, handler, filters...)
}

func (r RouterFunc) Options(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(OPTIONS, pattern, handler, filters...)
}

func (r RouterFunc) Head(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(HEAD, pattern, handler, filters...)
}

func (r RouterFunc) Any(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute("ANY", pattern, handler, filters...)
}

func (r RouterFunc) Setter() MethodSetter {
	return func(method Method) PathSetter {
		return func(pattern string) HandlerSetter {
			return func(handler interface{}, filters ...Filter) *swagger.Operation {
				return r.AddRoute(method, pattern, handler, filters...).SwaggerOperation()
			}
		}
	}
}

func (r RouterFunc) Group(pattern string, fn func(Router), filters ...Filter) {
	fn(RouterFunc(
		func(method Method, subpattern string, handler interface{}, subfilters ...Filter) *Leaf {
			return r.AddRoute(method, pattern+subpattern, handler, append(filters, subfilters...)...)
		}))
}

func (r RouterFunc) AddController(pattern string, controller ControllerInterface, filters ...Filter) {
	controllerFilter := &ControllerFilter{controller: controller}
	filters = append([]Filter{controllerFilter}, filters...)
	r.Group(pattern, func(sub Router) {
		for _, method := range Methods {
			handler := handlerForMethod(controller, method)
			leaf := Handle(sub, method, "/", handler)
			docMethodForMethod(controller, method)(leaf)
		}
	}, filters...)
}

func newRouter(n *Node) Router {
	return RouterFunc(func(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf {
		return n.AddRoute(method, pattern, handler, filters...)
	})
}
