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
	"sync"
)

// Context struct
type Context struct {
	Request         *http.Request
	Response        ResponseWriter
	Params          Params
	NotFoundHandler Handler
	data            map[string]interface{}
	Logger          Logger
	mutex           sync.RWMutex
}

func NewContext(w http.ResponseWriter, req *http.Request, logger Logger) *Context {
	return &Context{
		Request:  req,
		Response: NewResponseWriter(w),
		Params:   make(Params),
		data:     make(map[string]interface{}),
		Logger:   logger,
	}
}

func (ctx *Context) NotFound() {
	if ctx.NotFoundHandler == nil {
		http.NotFound(ctx.Response, ctx.Request)
		return
	}
	ctx.NotFoundHandler.Serve(ctx)
}

func (ctx *Context) Set(key string, value interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.data[key] = value
}

func (ctx *Context) Get(key string) interface{} {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	if v, ok := ctx.data[key]; ok {
		return v
	}
	return nil
}

func (ctx *Context) GetOK(key string) (value interface{}, ok bool) {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	if v, ok := ctx.data[key]; ok {
		return v, true
	}
	return nil, false
}

func (ctx *Context) Delete(key string) interface{} {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if v, ok := ctx.data[key]; ok {
		delete(ctx.data, key)
		return v
	}
	return nil
}
