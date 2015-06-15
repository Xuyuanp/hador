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
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	convey.Convey("Given a new Hodor h", t, func() {
		h := New()
		buf := make([]byte, 4096*4)
		buffer := bytes.NewBuffer(buf)
		logger := &logger{
			Logger: log.New(buffer, "[Hodor] ", 0),
		}
		convey.Convey("Should log messages", func() {
			h.Before(NewLogFilter(logger))
			h.Get("/500", HandlerFunc(func(ctx *Context) {
				ctx.WriteHeader(500)
			}))
			h.Get("/400", HandlerFunc(func(ctx *Context) {
				ctx.WriteHeader(400)
			}))
			h.Get("/300", HandlerFunc(func(ctx *Context) {
				ctx.WriteHeader(300)
			}))
			h.Get("/200", HandlerFunc(func(ctx *Context) {
				ctx.WriteHeader(200)
			}))
			h.Get("/100", HandlerFunc(func(ctx *Context) {
				ctx.WriteHeader(0)
			}))
			convey.Convey("Test log 500", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/500?a=1", nil)
				h.ServeHTTP(resp, req)
				convey.So(buffer.String(), convey.ShouldNotEqual, "")
			})
			convey.Convey("Test log 400", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/400?a=1", nil)
				h.ServeHTTP(resp, req)
				convey.So(buffer.String(), convey.ShouldNotEqual, "")
			})
			convey.Convey("Test log 300", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/300?a=1", nil)
				h.ServeHTTP(resp, req)
				convey.So(buffer.String(), convey.ShouldNotEqual, "")
			})
			convey.Convey("Test log 200", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/200?a=1", nil)
				h.ServeHTTP(resp, req)
				convey.So(buffer.String(), convey.ShouldNotEqual, "")
			})
			convey.Convey("Test log other", func() {
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/100?a=1", nil)
				h.ServeHTTP(resp, req)
				convey.So(buffer.String(), convey.ShouldNotEqual, "")
			})
		})
	})
}
