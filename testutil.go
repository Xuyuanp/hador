/*
 * Copyright 2014 Xuyuan Pang <xuyuanp # gmail dot com>
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

import "fmt"

func newSimpleHandler(content string) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.WriteString(content)
	})
}

func newParamHandler(name, def string) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx.WriteString(ctx.Params().GetStringMust(name, def))
	})
}

func echoHandler(ctx *Context) {
	ctx.WriteString(ctx.Request.RequestURI)
}

func emptyHandler(*Context) {}

func newMiddleware(name string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) {
			ctx.WriteString(fmt.Sprintf("%s Before -> ", name))
			next.Serve(ctx)
			ctx.WriteString(fmt.Sprintf(" -> %s After", name))
		})
	}
}
