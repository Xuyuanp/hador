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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestHodor(t *testing.T) {
	convey.Convey("Default ShouldNotBeNil", t, func() {
		convey.So(Default(), convey.ShouldNotBeNil)
	})
	convey.Convey("Given a new Hodor", t, func() {
		h := New()
		convey.Convey("h should not be nil", func() {
			convey.So(h, convey.ShouldNotBeNil)
		})
		convey.Convey("Run should not panic", func() {
			go h.Run(":6789")
		})
		convey.Convey("Test basic use", func() {
			h.Get("/hello", newSimpleHandler("hello"))
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hello", nil)
			h.ServeHTTP(resp, req)
			convey.So(resp.Body.String(), convey.ShouldEqual, "hello")
		})
		convey.Convey("Test Before method", func() {
			h.AddFilters(FilterFunc(func(ctx *Context, next Handler) {
				ctx.Response.Write([]byte("before"))
			}))
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			h.ServeHTTP(resp, req)
			convey.So(resp.Body.String(), convey.ShouldEqual, "before")
		})
		convey.Convey("Test bench", func() {
			h := New()
			h.Get("/hello/", newSimpleHandler("hello"))
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hello/", nil)
			h.ServeHTTP(resp, req)
			convey.So(resp.Body.String(), convey.ShouldEqual, "hello")
		})
		convey.Convey("Test swagger", func() {
			h.Get("/hello", newSimpleHandler("hello")).
				SwaggerOperation().
				DocSumDesc("hello", "get hello")
			h.SwaggerDocument().
				DocHost("127.0.0.1")
			h.Swagger(SwaggerConfig{APIPath: "/apidocs.json"})
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/apidocs.json", nil)
			h.ServeHTTP(resp, req)
			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
		})
		convey.Convey("Test Github", func() {
			h.Get("/repos/{owner}/{repo}/milestones/{number}", func(ctx *Context) {
				owner, _ := ctx.Params().GetString("owner")
				repo, _ := ctx.Params().GetString("repo")
				number, _ := ctx.Params().GetString("number")
				ctx.WriteString(fmt.Sprintf("%s:%s:%s", owner, repo, number))
			})
			h.Delete("/repos/{owner}/{repo}/milestones/{number}", func(ctx *Context) {
				owner, _ := ctx.Params().GetString("owner")
				repo, _ := ctx.Params().GetString("repo")
				number, _ := ctx.Params().GetString("number")
				ctx.WriteString(fmt.Sprintf("%s%s%s", owner, repo, number))
			})
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/repos/:owner/:repo/milestones/:number", nil)
			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(resp.Body.String(), convey.ShouldEqual, ":owner:repo:number")
		})
		convey.Convey("Test 405", func() {
			h := New()
			h.AddRoute(Method("GET"), "/fuck", func(ctx *Context) {
				ctx.WriteString(ctx.Request.URL.Path)
			})
			h.AddRoute(Method("DELETE"), "/fuck", func(ctx *Context) {
				ctx.WriteString(ctx.Request.URL.Path)
			})
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/fuck", nil)
			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(resp.Body.String(), convey.ShouldEqual, "/fuck")
		})
	})
}
