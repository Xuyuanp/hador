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

func TestSetter(t *testing.T) {
	convey.Convey("Test Setter", t, func() {
		h := New()
		h.Setter().
			Method(GET).
			Pattern("/hello/{name}").
			Handler(
			func(ctx *Context) {
				name, _ := ctx.Params().GetString("name")
				ctx.WriteString(name)
			})

		req, _ := http.NewRequest("GET", "/hello/jack", nil)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
		convey.So(resp.Body.String(), convey.ShouldEqual, "jack")
	})
}
