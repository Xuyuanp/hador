package hador

import "github.com/Xuyuanp/hador/swagger"

// HandlerSetter easy way to set Handler for a route.
type HandlerSetter func(handler interface{}, filters ...Filter) *swagger.Operation

// Handler calls HandlerSetter function.
func (hs HandlerSetter) Handler(handler interface{}, filters ...Filter) *swagger.Operation {
	return hs(handler)
}

// PathSetter easy way to set Path for a route.
type PathSetter func(path string) HandlerSetter

// Path calls PathSetter function.
func (ps PathSetter) Path(path string) HandlerSetter {
	return ps(path)
}

// MethodSetter easy way to set Method for a route.
type MethodSetter func(method Method) PathSetter

// Method calls MethodSetter function.
func (ms MethodSetter) Method(method Method) PathSetter {
	return ms(method)
}
