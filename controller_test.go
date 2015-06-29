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
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type testController struct {
	BaseController
	prepared bool
}

func (c *testController) Prepare(ctx *Context) bool {
	if !c.prepared {
		ctx.WriteString("not prepared")
	}
	return c.prepared
}

func (c *testController) Get(ctx *Context) {
	ctx.WriteHeader(http.StatusOK)
}

func TestController(t *testing.T) {
	convey.Convey("Test Controller", t, func() {
		convey.Convey("Test BaseController", func() {
			controller := &BaseController{}
			h := New()
			h.Any("/controller", controller)

			methods := map[string]int{
				"OPTIONS": http.StatusNotImplemented,
				"GET":     http.StatusNotImplemented,
				"HEAD":    http.StatusNotImplemented,
				"POST":    http.StatusNotImplemented,
				"PUT":     http.StatusNotImplemented,
				"DELETE":  http.StatusNotImplemented,
				"TRACE":   http.StatusNotImplemented,
				"CONNECT": http.StatusNotImplemented,
				"PATCH":   http.StatusNotImplemented,
			}

			for method, status := range methods {
				convey.Convey(fmt.Sprintf("Test %s", method), func() {
					req, _ := http.NewRequest(method, "/controller", nil)
					resp := httptest.NewRecorder()
					h.ServeHTTP(resp, req)
					convey.So(resp.Code, convey.ShouldEqual, status)
				})
			}
		})
		convey.Convey("Test testController", func() {
			convey.Convey("Test method override", func() {
				controller := &testController{prepared: true}
				h := New()
				h.Any("/controller", controller)
				req, _ := http.NewRequest("GET", "/controller", nil)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			})
			convey.Convey("Test not prepared", func() {
				controller := &testController{prepared: false}
				h := New()
				h.Any("/controller", controller)
				req, _ := http.NewRequest("GET", "/controller", nil)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)
				convey.So(resp.Body.String(), convey.ShouldEqual, "not prepared")
			})
			convey.Convey("Test illeagel method", func() {
				controller := &testController{prepared: true}
				h := New()
				h.Any("/controller", controller)
				req, _ := http.NewRequest("FOO", "/controller", nil)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusBadRequest)
				convey.So(resp.Body.String(), convey.ShouldEqual, "no such method: FOO\n")
			})
		})
	})
}
