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

	DocOptions(leaf *Leaf)
	DocGet(leaf *Leaf)
	DocHead(leaf *Leaf)
	DocPost(leaf *Leaf)
	DocPut(leaf *Leaf)
	DocDelete(leaf *Leaf)
	DocTrace(leaf *Leaf)
	DocConnect(leaf *Leaf)
	DocPatch(leaf *Leaf)
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

// DocOptions implements ControllerInterface.
func (c *BaseController) DocOptions(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocGet implements ControllerInterface.
func (c *BaseController) DocGet(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocHead implements ControllerInterface.
func (c *BaseController) DocHead(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocPost implements ControllerInterface.
func (c *BaseController) DocPost(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocPut implements ControllerInterface.
func (c *BaseController) DocPut(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocDelete implements ControllerInterface.
func (c *BaseController) DocDelete(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocTrace implements ControllerInterface.
func (c *BaseController) DocTrace(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocConnect implements ControllerInterface.
func (c *BaseController) DocConnect(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// DocPatch implements ControllerInterface.
func (c *BaseController) DocPatch(leaf *Leaf) {
	leaf.DocIgnore(true)
}

// ControllerFilter filters controller by using controller.Prepare methods.
type ControllerFilter struct {
	controller ControllerInterface
}

// Filter implements Filter interface.
func (cf *ControllerFilter) Filter(ctx *Context, next Handler) {
	if cf.controller.Prepare(ctx) {
		next.Serve(ctx)
	}
}

func methodForMethod(controller ControllerInterface, method Method) func(ctx *Context) {
	var fn func(ctx *Context)
	switch method {
	case OPTIONS:
		fn = controller.Options
	case GET:
		fn = controller.Get
	case HEAD:
		fn = controller.Head
	case POST:
		fn = controller.Post
	case PUT:
		fn = controller.Put
	case DELETE:
		fn = controller.Delete
	case TRACE:
		fn = controller.Trace
	case CONNECT:
		fn = controller.Connect
	case PATCH:
		fn = controller.Patch
	}
	return fn
}

func docMethodForMethod(controller ControllerInterface, method Method) func(leaf *Leaf) {
	var fn func(leaf *Leaf)
	switch method {
	case OPTIONS:
		fn = controller.DocOptions
	case GET:
		fn = controller.DocGet
	case HEAD:
		fn = controller.DocHead
	case POST:
		fn = controller.DocPost
	case PUT:
		fn = controller.DocPut
	case DELETE:
		fn = controller.DocDelete
	case TRACE:
		fn = controller.DocTrace
	case CONNECT:
		fn = controller.DocConnect
	case PATCH:
		fn = controller.DocPatch
	}
	return fn
}

func handlerForMethod(controller ControllerInterface, method Method) Handler {
	return HandlerFunc(methodForMethod(controller, method))
}
