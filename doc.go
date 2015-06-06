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

// Package hador is a high preformance and easy to use web framework in Go.
//
// example:
//
//	package main
//
//	import "github.com/Xuyuanp/hador"
//
//	func main() {
//		h := hador.New()
//
//		h.AddFilters(
//			h.NewLogFilter(h.Logger),
//		)
//
//		h.Group("/v1", func(root hador.Router) {
//
//			root.Get("/hello", hador.HandlerFunc(func(ctx hador.Context) {
//				ctx.Response.Write([]byte("hello"))
//			}), f2)
//
//			root.Get(`/hello/{name:\w+}`, hador.HandlerFunc(func(ctx hador.Context) {
//				name := ctx.Params().GetStringMust("name", "")
//				ctx.Response.Write([]byte("hello " + name))
//			}), f3, f4)
//
//		}, f1)
//
//		h.Get("/hello", hador.HandlerFunc(func(ctx hador.Context) {
//			ctx.Response.Write([]byte("hello"))
//		}), f5)
//
//		h.Run(":<your_port>")
//	}
//
// GET /hello
//
// 	LogFilter -> f5 -> handler -> f5 -> LogFilter
//
// GET /v1/hello
//
// 	LogFilter -> f1 -> f2 -> handler -> f2 -> f1 -> LogFilter
//
// GET /v1/hello/alice
//
// 	LogFilter -> f1 -> f3 -> f4 -> handler -> f4 -> f3 -> f1 -> LogFilter
//
//
package hador
