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
	"encoding/json"
	"encoding/xml"
	"fmt"
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

	errHandlers   map[int]func(...interface{})
	Err4XXHandler func(int, ...interface{})
	Err5XXHandler func(int, ...interface{})
}

func newContext(logger Logger) *Context {
	return &Context{
		Logger: logger,
	}
}

func (ctx *Context) reset(w ResponseWriter, req *http.Request) {
	ctx.Request = req
	ctx.Response = w
	ctx.params = nil
	ctx.data = nil
	ctx.errHandlers = nil
	ctx.Err4XXHandler = nil
	ctx.Err5XXHandler = nil
	ctx.segments = genSegments(req.URL.Path)
}

// OnError handles http error by calling handler registered in SetErrorHandler methods.
// If no handler registered with this status and noting written yet, http.Error would be used.
func (ctx *Context) OnError(status int, args ...interface{}) {
	// do nothing if not an error
	if status < 400 {
		return
	}
	// try to use custom error handler
	if ctx.errHandlers != nil {
		if h, ok := ctx.errHandlers[status]; ok {
			h(args...)
			return
		}
	}

	if status >= 400 && status < 500 && ctx.Err4XXHandler != nil {
		ctx.Err4XXHandler(status, args...)
		return
	}
	if status >= 500 && ctx.Err5XXHandler != nil {
		ctx.Err5XXHandler(status, args...)
		return
	}

	// use default http error
	switch status {
	case http.StatusMethodNotAllowed:
		// set Allow header for 405
		if len(args) > 0 {
			if allows, ok := args[0].([]string); ok && len(allows) > 0 {
				ctx.Response.Header().Set("Allow",
					strings.Join(allows, ","))
				args = args[1:]
			}
		}
	}
	if !ctx.Response.Written() {
		text := http.StatusText(status)
		if len(args) > 0 {
			if len(args) == 1 {
				text = fmt.Sprintf("%s", args[0])
			} else {
				text = fmt.Sprintf("%s", args)
			}
		}
		http.Error(ctx.Response, text, status)
	}
}

// SetErrorHandler sets custom handler for each http error
func (ctx *Context) SetErrorHandler(status int, handler func(...interface{})) {
	if ctx.errHandlers == nil {
		ctx.errHandlers = make(map[int]func(...interface{}))
	}
	ctx.errHandlers[status] = handler
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

// Write writes data into reponse by calling ctx.Response.Write method.
func (ctx *Context) Write(p []byte) (n int, err error) {
	return ctx.Response.Write(p)
}

// WriteStatus writes data into response by calling ctx.Response.Write method,
// and sets status as provided, if multi status provided, the first one will be used,
// if no status provided, does noting.
func (ctx *Context) WriteStatus(p []byte, status ...int) (n int, err error) {
	if status != nil && len(status) > 0 {
		ctx.Response.WriteHeader(status[0])
	}
	return ctx.Write(p)
}

// WriteString writes string into response by calling ctx.Write method.
func (ctx *Context) WriteString(s string) (n int, err error) {
	return ctx.Write([]byte(s))
}

// RenderJSON renders v in JSON format and sets status if provided.
func (ctx *Context) RenderJSON(v interface{}, status ...int) error {
	return ctx.renderJSON(v, false, status...)
}

// RenderPrettyJSON renders v in pretty JSON format and sets status if provided.
func (ctx *Context) RenderPrettyJSON(v interface{}, status ...int) error {
	return ctx.renderJSON(v, true, status...)
}

func (ctx *Context) renderJSON(v interface{}, indent bool, status ...int) error {
	charset := "utf-8"
	ctype := fmt.Sprintf("application/json; charset=%s", charset)
	return ctx.render(v, jsonMarshaler(indent), ctype, status...)
}

// RenderXML renders v in XML format and sets status if provided.
func (ctx *Context) RenderXML(v interface{}, status ...int) error {
	return ctx.renderXML(v, false, status...)
}

// RenderPrettyXML renders v in pretty XML format and sets status if provided.
func (ctx *Context) RenderPrettyXML(v interface{}, status ...int) error {
	return ctx.renderXML(v, true, status...)
}

func (ctx *Context) renderXML(v interface{}, indent bool, status ...int) error {
	charset := "utf-8"
	ctype := fmt.Sprintf("application/xml; charset=%s", charset)
	return ctx.render(v, xmlMarshaler(indent), ctype, status...)
}

type marshaler func(interface{}) ([]byte, error)

func jsonMarshaler(indent bool) marshaler {
	return func(v interface{}) ([]byte, error) {
		if indent {
			return json.MarshalIndent(v, "", "\t")
		}
		return json.Marshal(v)
	}
}

func xmlMarshaler(indent bool) marshaler {
	return func(v interface{}) ([]byte, error) {
		if indent {
			return xml.MarshalIndent(v, "", "\t")
		}
		return xml.Marshal(v)
	}
}

func (ctx *Context) render(v interface{}, m marshaler, ctype string, status ...int) error {
	data, err := m(v)
	if err != nil {
		return err
	}

	ctx.Response.Header().Set("Content-Type", ctype)
	_, err = ctx.WriteStatus(data, status...)
	return err
}

func (ctx *Context) ResolveJSON(v interface{}) error {
	return json.NewDecoder(ctx.Request.Body).Decode(v)
}

func (ctx *Context) ResolveXML(v interface{}) error {
	return xml.NewDecoder(ctx.Request.Body).Decode(v)
}
