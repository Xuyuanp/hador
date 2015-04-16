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

func TestNode(t *testing.T) {
	convey.Convey("Test tree", t, func() {
		convey.Convey("Test All method", func() {
			n := NewNode("", 0)
			n.Options("/a/b", newSimpleHandler("OPTIONS"))
			n.Get("/a/b", newSimpleHandler("GET"))
			n.Head("/a/b", newSimpleHandler("HEAD"))
			n.Post("/a/b", newSimpleHandler("POST"))
			n.Put("/a/b", newSimpleHandler("PUT"))
			n.Delete("/a/b", newSimpleHandler("DELETE"))
			n.Trace("/a/b", newSimpleHandler("TRACE"))
			n.Connect("/a/b", newSimpleHandler("CONNECT"))
			n.Patch("/a/b", newSimpleHandler("PATCH"))
			n.Any("/a/c", newSimpleHandler("ANY"))
			convey.Convey("Test OPTIONS", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("OPTIONS", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "OPTIONS")
			})
			convey.Convey("Test Get", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "GET")
			})
			convey.Convey("Test HEAD", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("HEAD", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "HEAD")
			})
			convey.Convey("Test POST", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "POST")
			})
			convey.Convey("Test PUT", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("PUT", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "PUT")
			})
			convey.Convey("Test DELETE", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("DELETE", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "DELETE")
			})
			convey.Convey("Test TRACE", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("TRACE", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "TRACE")
			})
			convey.Convey("Test CONNECT", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("CONNECT", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "CONNECT")
			})
			convey.Convey("Test PATCH", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("PATCH", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "PATCH")
			})
			convey.Convey("Test ANY", func() {
				convey.Convey("GET", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/a/c", nil)
					ctx := NewContext(resp, req, defaultLogger)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "ANY")
				})
				convey.Convey("POST", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("POST", "/a/c", nil)
					ctx := NewContext(resp, req, defaultLogger)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "ANY")
				})
			})
			convey.Convey("Test Group", func() {
				n := NewNode("", 0)
				n.Group("/a", func(r Router) {
					r.Get("/b", newSimpleHandler("GET"))
					r.Post("/c", newSimpleHandler("POST"))
				})
				convey.Convey("/a/b", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/a/b", nil)
					ctx := NewContext(resp, req, defaultLogger)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "GET")
				})
				convey.Convey("/a/c", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("POST", "/a/c", nil)
					ctx := NewContext(resp, req, defaultLogger)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "POST")
				})
			})
			convey.Convey("Register multi handlers with same method and same path would case panic", func() {
				defer func() {
					convey.So(recover(), convey.ShouldNotBeNil)
				}()
				n := NewNode("", 0)
				n.Get("/test", newSimpleHandler("h1"))
				n.Get("/test", newSimpleHandler("h2"))
			})
		})
		convey.Convey("Test regexp path", func() {
			n := NewNode("", 0)
			n.Get(`/(?P<name>\w+)`, newSimpleHandler("h1"))
			n.Get(`/(?P<name>\w+)/(?P<age>[1-9]\d*)`, newSimpleHandler("h2"))
			convey.Convey("/jack", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/jack", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "h1")
			})
			convey.Convey("/jack/12", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/jack/12", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "h2")
			})
		})
		convey.Convey("Test Filter", func() {
			n := NewNode("", 0)
			n.Get("/a", newSimpleHandler("h1"))
			n.Get("/a/b", newSimpleHandler("h2"), FilterFunc(func(ctx *Context, next Handler) {
				ctx.Response.Write([]byte("Filter"))
			}))
			convey.Convey("test /a", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/a", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "h1")
			})
			convey.Convey("test /a/b", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "Filter")
			})
		})
		convey.Convey("Test FilteFilter with Group", func() {
			n := NewNode("", 0)
			n.Group("/a", func(r Router) {
				r.Get("/b", newSimpleHandler("h1"))
				r.Get("/c", newSimpleHandler("h2"))
			}, FilterFunc(func(ctx *Context, next Handler) {
				ctx.Response.Write([]byte("Filter"))
			}))
			convey.Convey("test /a/b", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/a/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "Filter")
			})
			convey.Convey("test /a/c", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/a/c", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "Filter")
			})
		})
		convey.Convey("Test error", func() {
			n := NewNode("", 0)
			n.Get("/a", newSimpleHandler("GET"))
			convey.Convey("test GET /a", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/a", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "GET")
			})
			convey.Convey("test POST /a", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/a", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusMethodNotAllowed)
			})
			convey.Convey("test GET /b", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/b", nil)
				ctx := NewContext(resp, req, defaultLogger)
				n.Serve(ctx)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
			})
		})
		convey.Convey("Test Path", func() {
			n := NewNode("", 0)
			l := n.Get("/", newSimpleHandler("GET"))
			l1 := n.Get("/a/b/c/d", newSimpleHandler("GET"))
			l2 := n.Get("/a/(?P<name>.+)", newSimpleHandler("GET"))
			convey.So(l.Path(), convey.ShouldEqual, "/")
			convey.So(l1.Path(), convey.ShouldEqual, "/a/b/c/d")
			convey.So(l2.Path(), convey.ShouldEqual, "/a/(?P<name>.+)")
		})
	})
}
