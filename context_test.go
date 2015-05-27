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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {
	convey.Convey("given a new context", t, func() {
		ctx := newContext(defaultLogger)
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		convey.So(err, convey.ShouldBeNil)
		rw := NewResponseWriter(resp)
		ctx.reset(rw, req)

		convey.Convey("test internal data", func() {
			convey.Convey("get & delete should be nil", func() {
				convey.So(ctx.Get("foo"), convey.ShouldBeNil)
				_, ok := ctx.GetOK("foo")
				convey.So(ok, convey.ShouldBeFalse)
				convey.So(ctx.Delete("foo"), convey.ShouldBeNil)
				convey.Convey("set key-value", func() {
					ctx.Set("foo", "bar")
					convey.Convey("get & delete should be ok", func() {
						convey.So(ctx.Get("foo"), convey.ShouldEqual, "bar")
						_, ok := ctx.GetOK("foo")
						convey.So(ok, convey.ShouldBeTrue)

						convey.So(ctx.Delete("foo"), convey.ShouldEqual, "bar")
					})
					convey.Convey("get & delete unknown key should be nil", func() {
						convey.So(ctx.Get("bazz"), convey.ShouldBeNil)
						_, ok := ctx.GetOK("bazz")
						convey.So(ok, convey.ShouldBeFalse)

						convey.So(ctx.Delete("bazz"), convey.ShouldBeNil)
					})
				})
			})
		})

		convey.Convey("test error handler", func() {
			convey.Convey("test custome error handler", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)
				ctx.SetErrorHandler(http.StatusNotFound, HandlerFunc(func(c *Context) {
					c.Response.WriteHeader(http.StatusNotFound)
					c.Response.Write([]byte("404"))
				}))
				ctx.OnError(http.StatusNotFound)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
				convey.So(resp.Body.String(), convey.ShouldEqual, "404")
			})
			convey.Convey("test default error handler", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)
				ctx.OnError(http.StatusNotFound)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
				convey.So(resp.Body.String(), convey.ShouldEqual,
					http.StatusText(http.StatusNotFound)+"\n")
			})
		})
	})
}
