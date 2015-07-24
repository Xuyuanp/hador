package hador

// HandlerSetter easy way to set Handler for a route.
type HandlerSetter func(handler interface{}, filters ...Filter) *Leaf

// Handler calls HandlerSetter function.
func (hs HandlerSetter) Handler(handler interface{}, filters ...Filter) *Leaf {
	return hs(handler)
}

// PatternSetter easy way to set Path for a route.
type PatternSetter func(pattern string) HandlerSetter

// Pattern calls Pattern function.
func (ps PatternSetter) Pattern(pattern string) HandlerSetter {
	return ps(pattern)
}

// MethodSetter easy way to set Method for a route.
type MethodSetter func(method Method) PatternSetter

// Method calls MethodSetter function.
func (ms MethodSetter) Method(method Method) PatternSetter {
	return ms(method)
}
