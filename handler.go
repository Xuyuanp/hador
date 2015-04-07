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

import "strings"

type Handler interface {
	Serve(*Context)
}

type HandlerFunc func(ctx *Context)

func (hf HandlerFunc) Serve(ctx *Context) {
	hf(ctx)
}

type MethodHandler struct {
	handlers map[string]Handler
}

func NewMethodHandler() *MethodHandler {
	return &MethodHandler{
		handlers: make(map[string]Handler),
	}
}

func (h *MethodHandler) Serve(ctx *Context) {
	if h.IsEmpty() {
		ctx.NotFound()
	} else if handler, ok := h.handlers[ctx.Request.Method]; ok {
		handler.Serve(ctx)
	} else if handler, ok := h.handlers["ANY"]; ok {
		handler.Serve(ctx)
	} else {
		h.MethodNotAllowed(ctx)
	}
}

func (h *MethodHandler) IsEmpty() bool {
	return len(h.handlers) == 0
}

func (h *MethodHandler) Handle(method string, handler Handler) bool {
	if _, ok := h.handlers[method]; ok {
		return false
	}
	h.handlers[method] = handler
	return true
}

func (h *MethodHandler) MethodNotAllowed(ctx *Context) {
	methods := []string{}
	for m := range h.handlers {
		methods = append(methods, m)
	}
	allow := strings.Join(methods, ",")
	ctx.MethodNotAllowed(allow)
}
