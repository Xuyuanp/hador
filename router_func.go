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

var _ Router = RouterFunc(nil)

// RouterFunc is a function type implemented Router interface.
type RouterFunc func(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf

// AddRoute calls RouterFunc function. It is the most important method of RouterFunc.
// All other methods call this method finally.
func (r RouterFunc) AddRoute(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r(method, pattern, handler, filters...)
}

// Route returns a setter-chain to add a new route step-by-step.
func (r RouterFunc) Route() MethodSetter {
	return func(method Method) PatternSetter {
		return func(pattern string) HandlerSetter {
			return func(handler interface{}, filters ...Filter) *Leaf {
				return r.AddRoute(method, pattern, handler, filters...)
			}
		}
	}
}

// Get adds a new route binded with GET method.
func (r RouterFunc) Get(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Get().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Post adds a new route binded with POST method.
func (r RouterFunc) Post(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Post().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Put adds a new route binded with PUT method.
func (r RouterFunc) Put(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Put().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Delete adds a new route binded with DELETE method.
func (r RouterFunc) Delete(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Delete().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Patch adds a new route binded with Patch method.
func (r RouterFunc) Patch(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Patch().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Trace adds a new route binded with TRACE method.
func (r RouterFunc) Trace(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Trace().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Connect adds a new route binded with CONNECT method.
func (r RouterFunc) Connect(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Connect().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Options adds a new route binded with OPTIONS method.
func (r RouterFunc) Options(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Options().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Head adds a new route binded with HEAD method.
func (r RouterFunc) Head(pattern string, handler interface{}, filters ...Filter) *Leaf {
	return r.Route().Head().Pattern(pattern).Handler(handler).AddFilters(filters...)
}

// Any adds a new route binded with all method.
func (r RouterFunc) Any(pattern string, handler interface{}, filters ...Filter) *Leaf {
	for _, method := range Methods {
		r.Route().Method(method).Pattern(pattern).Handler(handler).AddFilters(filters...)
	}
	return nil
}

// Group adds multi routes one time.
func (r RouterFunc) Group(pattern string, fn func(Router), filters ...Filter) {
	fn(RouterFunc(
		func(method Method, subpattern string, handler interface{}, subfilters ...Filter) *Leaf {
			return r.Route().
				Method(method).
				Pattern(pattern + subpattern).
				Handler(handler).
				AddFilters(filters...).
				AddFilters(subfilters...)
		}))
}

// AddController adds routes of all methods by calling controller's matched method.
func (r RouterFunc) AddController(pattern string, controller ControllerInterface, filters ...Filter) {
	controllerFilter := &ControllerFilter{controller: controller}
	filters = append([]Filter{controllerFilter}, filters...)
	r.Group(pattern, func(sub Router) {
		for _, method := range Methods {
			handler := handlerForMethod(controller, method)
			leaf := sub.Route().Method(method).Pattern("/").Handler(handler)
			docFn := docMethodForMethod(controller, method)
			docFn(leaf)
		}
	}, filters...)
}
