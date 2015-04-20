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
	Request       *http.Request
	Response      ResponseWriter
	Params        Params
	ErrorHandlers map[int]Handler
	data          map[string]interface{}
	segments      []string
	Logger        Logger
}

// NewContext creates new Context instance
func NewContext(w http.ResponseWriter, req *http.Request, logger Logger) *Context {
	return &Context{
		Request:       req,
		Response:      NewResponseWriter(w),
		Params:        make(Params),
		ErrorHandlers: make(map[int]Handler),
		data:          make(map[string]interface{}),
		segments:      genSegments(req.URL.Path),
		Logger:        logger,
	}
}

// NotFound handles 404 error
func (ctx *Context) NotFound() {
	if h, ok := ctx.ErrorHandlers[http.StatusNotFound]; ok {
		h.Serve(ctx)
		return
	}
	http.NotFound(ctx.Response, ctx.Request)
}

// MethodNotAllowed handles 405 error
func (ctx *Context) MethodNotAllowed(allow []string) {
	ctx.Response.Header().Set("Allow", strings.Join(allow, ","))
	if h, ok := ctx.ErrorHandlers[http.StatusMethodNotAllowed]; ok {
		h.Serve(ctx)
		return
	}
	http.Error(ctx.Response, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// Set saves data in the context
func (ctx *Context) Set(key string, value interface{}) {
	ctx.data[key] = value
}

// Get retrieves data from the context
func (ctx *Context) Get(key string) interface{} {
	if v, ok := ctx.data[key]; ok {
		return v
	}
	return nil
}

// GetOK retrieves data from the context, and returns (nil, false) if no data
func (ctx *Context) GetOK(key string) (value interface{}, ok bool) {
	if v, ok := ctx.data[key]; ok {
		return v, true
	}
	return nil, false
}

// Delete removes data from the context
func (ctx *Context) Delete(key string) interface{} {
	if v, ok := ctx.data[key]; ok {
		delete(ctx.data, key)
		return v
	}
	return nil
}
