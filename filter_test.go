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

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestFilterChain(t *testing.T) {
	convey.Convey("Given Filter f1, f2 and Handler h", t, func() {
		f1 := FilterFunc(func(ctx *Context, next Handler) {
			ctx.Response.Write([]byte("f1Before"))
			next.Serve(ctx)
			ctx.Response.Write([]byte("f1After"))
		})
		f2 := FilterFunc(func(ctx *Context, next Handler) {
			ctx.Response.Write([]byte("f2Before"))
			next.Serve(ctx)
			ctx.Response.Write([]byte("f2After"))
		})
		h := HandlerFunc(func(ctx *Context) {
			ctx.Response.Write([]byte("handler"))
		})
		convey.Convey("Test panic", func() {
			convey.Convey("NewFilterChain with nil Handler should panic", func() {
				defer func() {
					convey.So(recover(), convey.ShouldNotBeNil)
				}()
				NewFilterChain(nil)
			})
		})
		convey.Convey("Test Insert", func() {
			convey.Convey("test InsertFront", func() {
				fc := NewFilterChain(h, f1)
				fc.InsertFront(f2)
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/", nil)
				ctx := newContext(defaultLogger)
				ctx.reset(NewResponseWriter(resp), req)
				fc.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "f2Beforef1Beforehandlerf1Afterf2After")
			})
			convey.Convey("test InsertBack", func() {
				fc := NewFilterChain(h, f1)
				fc.InsertBack(f2)
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/", nil)
				ctx := newContext(defaultLogger)
				ctx.reset(NewResponseWriter(resp), req)
				fc.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "f1Beforef2Beforehandlerf2Afterf1After")
			})
		})
		convey.Convey("Test FilterChain without Filter", func() {
			fc := NewFilterChain(h)
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			ctx := newContext(defaultLogger)
			ctx.reset(NewResponseWriter(resp), req)
			fc.Serve(ctx)
			convey.Convey("response string should be handler", func() {
				convey.So(resp.Body.String(), convey.ShouldEqual, "handler")
			})
		})
		convey.Convey("NewFilterChain should not be nil", func() {
			fc := NewFilterChain(h)
			convey.So(fc, convey.ShouldNotBeNil)
			convey.Convey("Test Before", func() {
				fc.Before(f1)
				fc.BeforeFunc(f2)

				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/", nil)
				ctx := newContext(defaultLogger)
				ctx.reset(NewResponseWriter(resp), req)
				fc.Serve(ctx)
				convey.Convey("response string should be f1f2handler", func() {
					convey.So(resp.Body.String(), convey.ShouldEqual, "f1Beforef2Beforehandlerf2Afterf1After")
				})
			})
		})

		convey.Convey("Test CombineFilters", func() {
			fs := CombineFilters(f1, f2, nil)
			convey.Convey("Test Before", func() {

				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/", nil)
				ctx := newContext(defaultLogger)
				ctx.reset(NewResponseWriter(resp), req)
				fs.Filter(ctx, h)
				convey.Convey("response string should be f1f2handler", func() {
					convey.So(resp.Body.String(), convey.ShouldEqual, "f1Beforef2Beforehandlerf2Afterf1After")
				})
			})
		})
	})
}
