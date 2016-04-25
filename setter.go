/*
 * Copyright 2016 Xuyuan Pang
 * Author: Xuyuan Pang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hador

// HandlerSetter easy way to set Handler for a route.
type HandlerSetter func(handler interface{}, filters ...Filter) *Leaf

// Handler calls HandlerSetter function.
func (hs HandlerSetter) Handler(handler interface{}, filters ...Filter) *Leaf {
	return hs(handler)
}

// Filters creates a new HandlerSetter to append filters.
func (hs HandlerSetter) Filters(filters ...Filter) HandlerSetter {
	return func(handler interface{}, fs ...Filter) *Leaf {
		return hs(handler, append(filters, fs...)...)
	}
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

// Grouper type is a group routing tool.
type Grouper func(func(Router), ...Filter)

// For calls Grouper function.
func (g Grouper) For(fn func(Router)) {
	g(fn)
}

// Filters creates a new Grouper appending filters.
func (g Grouper) Filters(fs ...Filter) Grouper {
	return func(fn func(Router), filters ...Filter) {
		g(fn, append(fs, filters...)...)
	}
}

// Group creates a new Grouper with root.
func (ms MethodSetter) Group(root string) Grouper {
	return func(fn func(Router), fs ...Filter) {
		fn(RouterFunc(
			func(method Method, subpattern string, handler interface{}, subfilters ...Filter) *Leaf {
				return ms.Method(method).
					Pattern(root + subpattern).
					Filters(fs...).
					Filters(subfilters...).
					Handler(handler)
			}))
	}
}
