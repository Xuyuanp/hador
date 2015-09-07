package hador

var _ Router = RouterFunc(nil)

// RouterFunc is a function type implemented Router interface.
type RouterFunc func(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf

// AddRoute calls RouterFunc function. It is the most important method of RouterFunc.
// All other methods call this method finally.
func (r RouterFunc) AddRoute(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r(method, pattern, handler, filters...)
}

// Get adds a new route binded with GET method.
func (r RouterFunc) Get(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(GET, pattern, handler, filters...)
}

// Post adds a new route binded with POST method.
func (r RouterFunc) Post(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(POST, pattern, handler, filters...)
}

// Put adds a new route binded with PUT method.
func (r RouterFunc) Put(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(PUT, pattern, handler, filters...)
}

// Delete adds a new route binded with DELETE method.
func (r RouterFunc) Delete(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(DELETE, pattern, handler, filters...)
}

// Patch adds a new route binded with Patch method.
func (r RouterFunc) Patch(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(PATCH, pattern, handler, filters...)
}

// Trace adds a new route binded with TRACE method.
func (r RouterFunc) Trace(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(TRACE, pattern, handler, filters...)
}

// Connect adds a new route binded with CONNECT method.
func (r RouterFunc) Connect(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(CONNECT, pattern, handler, filters...)
}

// Options adds a new route binded with OPTIONS method.
func (r RouterFunc) Options(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(OPTIONS, pattern, handler, filters...)
}

// Head adds a new route binded with HEAD method.
func (r RouterFunc) Head(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.AddRoute(HEAD, pattern, handler, filters...)
}

// Any adds a new route binded with all method.
func (r RouterFunc) Any(pattern string, handler interface{}, filters ...Filter) *Leaf {
	for _, method := range Methods {
		r.AddRoute(method, pattern, handler, filters...)
	}
	return nil
}

// Setter returns a setter-chain to add a new route step-by-step.
func (r RouterFunc) Setter() MethodSetter {
	return func(method Method) PatternSetter {
		return func(pattern string) HandlerSetter {
			return func(handler interface{}, filters ...Filter) *Leaf {
				return r.AddRoute(method, pattern, handler, filters...)
			}
		}
	}
}

// Group adds multi routes one time.
func (r RouterFunc) Group(pattern string, fn func(Router), filters ...Filter) {
	fn(RouterFunc(
		func(method Method, subpattern string, handler interface{}, subfilters ...Filter) *Leaf {
			return r.AddRoute(method,
				pattern+subpattern,
				handler,
				append(filters, subfilters...)...)
		}))
}

// AddController adds routes of all methods by calling controller's matched method.
func (r RouterFunc) AddController(pattern string, controller ControllerInterface, filters ...Filter) {
	controllerFilter := &ControllerFilter{controller: controller}
	filters = append([]Filter{controllerFilter}, filters...)
	r.Group(pattern, func(sub Router) {
		for _, method := range Methods {
			handler := handlerForMethod(controller, method)
			leaf := sub.AddRoute(method, "/", handler)
			docFn := docMethodForMethod(controller, method)
			docFn(leaf)
		}
	}, filters...)
}
