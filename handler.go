/*
 * Copyright 2015 <name of copyright holder>
 * Author: Xuyuan Pang <xuyuanp # gmail dot com>
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

import "net/http"

// Handler interface
type Handler interface {
	Serve(*Context)
}

// HandlerFunc is Handle as function
type HandlerFunc func(ctx *Context)

// Serve implements Handler interface by calling HandlerFunc function
func (hf HandlerFunc) Serve(ctx *Context) {
	hf(ctx)
}

// Wrap wraps http.Handler to Handler
func Wrap(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		handler.ServeHTTP(ctx.Response, ctx.Request)
	}
}

// WrapFunc wraps http.HandlerFunc to HandlerFunc
func WrapFunc(hf func(http.ResponseWriter, *http.Request)) HandlerFunc {
	return Wrap(http.HandlerFunc(hf))
}

func parseHandler(h interface{}) Handler {
	switch v := h.(type) {
	case Handler:
		return v
	case func(*Context):
		return HandlerFunc(v)
	case http.Handler:
		return Wrap(v)
	case func(http.ResponseWriter, *http.Request):
		return WrapFunc(v)
	case ControllerInterface:
		return &ControllerHandler{v}
	}
	panic("invalid handler")
}
