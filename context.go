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

import (
	"net/http"
	"strings"
)

// Context struct
type Context struct {
	Request  *http.Request
	Response ResponseWriter
	params   Params

	data     map[string]interface{}
	segments []string
	Logger   Logger

	NotFoundHandler         Handler
	MethodNotAllowedHandler func(ctx *Context, allows []string)
}

func newContext(logger Logger) *Context {
	return &Context{
		Logger: logger,

		NotFoundHandler: HandlerFunc(func(ctx *Context) {
			http.NotFound(ctx.Response, ctx.Request)
		}),
		MethodNotAllowedHandler: func(ctx *Context, allows []string) {
			ctx.Response.Header().Set("Allow", strings.Join(allows, ","))
			http.Error(ctx.Response, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		},
	}
}

func (ctx *Context) reset(w ResponseWriter, req *http.Request) {
	ctx.Request = req
	ctx.Response = w
	ctx.params = nil
	ctx.data = nil
	ctx.segments = genSegments(req.URL.Path)

	ctx.NotFoundHandler = HandlerFunc(func(ctx *Context) {
		http.NotFound(ctx.Response, ctx.Request)
	})
	ctx.MethodNotAllowedHandler = func(ctx *Context, allows []string) {
		ctx.Response.Header().Set("Allow", strings.Join(allows, ","))
		http.Error(ctx.Response, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// NotFound handles 404 error
func (ctx *Context) NotFound() {
	ctx.NotFoundHandler.Serve(ctx)
}

// MethodNotAllowed handles 405 error
func (ctx *Context) MethodNotAllowed(allows []string) {
	ctx.MethodNotAllowedHandler(ctx, allows)
}

// Params returns params lazy-init
func (ctx *Context) Params() Params {
	if ctx.params == nil {
		ctx.params = make(Params)
	}
	return ctx.params
}

// Set saves data in the context
func (ctx *Context) Set(key string, value interface{}) {
	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	ctx.data[key] = value
}

// Get retrieves data from the context
func (ctx *Context) Get(key string) interface{} {
	if ctx.data == nil {
		return nil
	}
	if v, ok := ctx.data[key]; ok {
		return v
	}
	return nil
}

// GetOK retrieves data from the context, and returns (nil, false) if no data
func (ctx *Context) GetOK(key string) (value interface{}, ok bool) {
	if ctx.data == nil {
		return nil, false
	}
	if v, ok := ctx.data[key]; ok {
		return v, true
	}
	return nil, false
}

// Delete removes data from the context
func (ctx *Context) Delete(key string) interface{} {
	if ctx.data == nil {
		return nil
	}
	if v, ok := ctx.data[key]; ok {
		delete(ctx.data, key)
		return v
	}
	return nil
}
