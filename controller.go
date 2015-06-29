/*
 * Copyright 2015 Xuyuan Pang
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

import "net/http"

// ControllerInterface is definition of controller API.
type ControllerInterface interface {
	Prepare(ctx *Context) bool

	Options(ctx *Context)
	Get(ctx *Context)
	Head(ctx *Context)
	Post(ctx *Context)
	Put(ctx *Context)
	Delete(ctx *Context)
	Trace(ctx *Context)
	Connect(ctx *Context)
	Patch(ctx *Context)
}

// BaseController provides an empty controller,
// returns Not Implemented status in any method.
type BaseController struct {
}

// Prepare runs before other methods, and ends this request if returns false.
func (c *BaseController) Prepare(ctx *Context) bool {
	return true
}

// Options implements ControllerInterface.
func (c *BaseController) Options(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Get implements ControllerInterface.
func (c *BaseController) Get(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Head implements ControllerInterface.
func (c *BaseController) Head(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Post implements ControllerInterface.
func (c *BaseController) Post(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Put implements ControllerInterface.
func (c *BaseController) Put(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Delete implements ControllerInterface.
func (c *BaseController) Delete(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Trace implements ControllerInterface.
func (c *BaseController) Trace(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Connect implements ControllerInterface.
func (c *BaseController) Connect(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// Patch implements ControllerInterface.
func (c *BaseController) Patch(ctx *Context) {
	ctx.OnError(http.StatusNotImplemented)
}

// ControllerHandler is a Handler that dispatchers request by method into
// the innerController.
type ControllerHandler struct {
	innerController ControllerInterface
}

// Serve implements Handler interface.
func (ch *ControllerHandler) Serve(ctx *Context) {
	if !ch.innerController.Prepare(ctx) {
		return
	}
	switch method := ctx.Request.Method; method {
	case "OPTIONS":
		ch.innerController.Options(ctx)
	case "GET":
		ch.innerController.Get(ctx)
	case "HEAD":
		ch.innerController.Head(ctx)
	case "POST":
		ch.innerController.Post(ctx)
	case "PUT":
		ch.innerController.Put(ctx)
	case "DELETE":
		ch.innerController.Delete(ctx)
	case "TRACE":
		ch.innerController.Trace(ctx)
	case "CONNECT":
		ch.innerController.Connect(ctx)
	case "PATCH":
		ch.innerController.Patch(ctx)
	default:
		ctx.OnError(http.StatusBadRequest, "no such method: "+method)
	}
}
