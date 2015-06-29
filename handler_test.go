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

func TestHandler(t *testing.T) {
	convey.Convey("Test Handler", t, func() {
		convey.Convey("Test parseHandler", func() {
			parseHandler(func(ctx *Context) {})
			parseHandler(HandlerFunc(func(ctx *Context) {}))
			parseHandler(func(http.ResponseWriter, *http.Request) {})
			parseHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
			func() {
				defer func() {
					convey.So(recover(), convey.ShouldEqual, "invalid handler")
				}()
				parseHandler(1)
			}()
		})
		convey.Convey("Test Wrap", func() {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			}
			h := New()
			h.Get("/wrap", handler)

			req, _ := http.NewRequest("GET", "/wrap", nil)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			convey.So(resp.Body.String(), convey.ShouldEqual, "OK")
		})
	})
}
