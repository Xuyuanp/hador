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

func TestMiddleware(t *testing.T) {
	convey.Convey("Test Middleware", t, func() {
		m1 := newMiddleware("m1")
		handler := newSimpleHandler("handler")

		handler = m1(handler)

		h := New()

		h.Get("/", handler)

		req, _ := http.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)

		convey.So(resp.Body.String(), convey.ShouldEqual, "m1 Before -> handler -> m1 After")
	})

	convey.Convey("Test Middleware as Filter", t, func() {
		m1 := newMiddleware("m1")
		handler := newSimpleHandler("handler")

		h := New()

		h.Get("/", handler, m1)

		req, _ := http.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)

		convey.So(resp.Body.String(), convey.ShouldEqual, "m1 Before -> handler -> m1 After")
	})
}
