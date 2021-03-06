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

import (
	"fmt"
	"net/http"
)

// OK default handler
func OK(msg string) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.WriteString(msg, http.StatusOK)
	})
}

// Status handler
func Status(status int) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.WriteHeader(status)
	})
}

// JSON handler
func JSON(v interface{}) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.RenderJSON(v)
	})
}

// XML handler
func XML(v interface{}) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.RenderXML(v)
	})
}

// Data handler
func Data(v interface{}) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.WriteString(fmt.Sprintf("%v", v))
	})
}
