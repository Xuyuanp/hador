/*
 * Copyright 2015 Xuyuan Pang
 * Author: Pang Xuyuan <xuyuanp # gmail dot com>
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

// Method is HTTP method.
type Method string

// available methods.
const (
	OPTIONS Method = "OPTIONS"
	GET            = "GET"
	HEAD           = "HEAD"
	POST           = "POST"
	PUT            = "PUT"
	DELETE         = "DELETE"
	TRACE          = "TRACE"
	CONNECT        = "CONNECT"
	PATCH          = "PATCH"
)

// Methods is a list of all valid methods.
var Methods = []Method{
	OPTIONS,
	GET,
	HEAD,
	POST,
	PUT,
	DELETE,
	TRACE,
	CONNECT,
	PATCH,
}

func (m Method) String() string {
	return string(m)
}

// Router interface
type Router interface {
	Route() MethodSetter

	Options(pattern string, h interface{}, filters ...Filter) *Leaf
	Get(pattern string, h interface{}, filters ...Filter) *Leaf
	Head(pattern string, h interface{}, filters ...Filter) *Leaf
	Post(pattern string, h interface{}, filters ...Filter) *Leaf
	Put(pattern string, h interface{}, filters ...Filter) *Leaf
	Delete(pattern string, h interface{}, filters ...Filter) *Leaf
	Trace(pattern string, h interface{}, filters ...Filter) *Leaf
	Connect(pattern string, h interface{}, filters ...Filter) *Leaf
	Patch(pattern string, h interface{}, filters ...Filter) *Leaf

	// Any routes doesn't support swagger API, always returns nil.
	Any(pattern string, h interface{}, filters ...Filter) *Leaf

	// Be carefull! All requests match this pattern would be processed by the provided
	// controller instance, no new ControllerInterface implemented instance would be
	// created. If you want to share data in your controller, please make sure it is
	// goroutine-safe.
	AddController(pattern string, controller ControllerInterface, filters ...Filter)

	AddRoute(method Method, pattern string, h interface{}, filters ...Filter) *Leaf

	Group(pattern string, fn func(Router), filters ...Filter)
}
