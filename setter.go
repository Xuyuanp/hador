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

// Options short for Method(OPTIONS)
func (ms MethodSetter) Options() PatternSetter {
	return ms.Method(OPTIONS)
}

// Get short for Method(GET)
func (ms MethodSetter) Get() PatternSetter {
	return ms.Method(GET)
}

// Head short for Method(HEAD)
func (ms MethodSetter) Head() PatternSetter {
	return ms.Method(HEAD)
}

// Post short for Method(POST)
func (ms MethodSetter) Post() PatternSetter {
	return ms.Method(POST)
}

// Put short for Method(PUT)
func (ms MethodSetter) Put() PatternSetter {
	return ms.Method(PUT)
}

// Delete short for Method(DELETE)
func (ms MethodSetter) Delete() PatternSetter {
	return ms.Method(DELETE)
}

// Trace short for Method(TRACE)
func (ms MethodSetter) Trace() PatternSetter {
	return ms.Method(TRACE)
}

// Connect short for Method(CONNECT)
func (ms MethodSetter) Connect() PatternSetter {
	return ms.Method(CONNECT)
}

// Patch short for Method(PATCH)
func (ms MethodSetter) Patch() PatternSetter {
	return ms.Method(PATCH)
}
