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

func TestStatic(t *testing.T) {
	convey.Convey("test static", t, func() {
		convey.Convey("test static GET", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://127.0.0.1:4000/hador.go", nil)
			convey.So(err, convey.ShouldBeNil)
			h := New()
			s := NewStatic(http.Dir("."))
			h.Before(s)

			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(resp.Header().Get("Expires"), convey.ShouldBeBlank)
			convey.So(resp.Body.Len() > 0, convey.ShouldBeTrue)
		})
		convey.Convey("test static HEAD", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("HEAD", "http://127.0.0.1:4000/hador.go", nil)
			convey.So(err, convey.ShouldBeNil)
			h := New()
			s := NewStatic(http.Dir("."))
			h.Before(s)

			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(resp.Body.Len() == 0, convey.ShouldBeTrue)
		})
		convey.Convey("test static POST", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://127.0.0.1:4000/hador.go", nil)
			convey.So(err, convey.ShouldBeNil)
			h := New()
			s := NewStatic(http.Dir("."))
			h.Before(s)

			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
		})
		convey.Convey("test static index file", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://127.0.0.1:4000/", nil)
			convey.So(err, convey.ShouldBeNil)
			h := New()
			s := NewStatic(http.Dir("."))
			s.IndexFile = "hador.go"
			h.Before(s)

			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(resp.Header().Get("Expires"), convey.ShouldBeBlank)
			convey.So(resp.Body.Len() > 0, convey.ShouldBeTrue)
		})
		convey.Convey("test static prefix", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://127.0.0.1:4000/public/hador.go", nil)
			convey.So(err, convey.ShouldBeNil)
			h := New()
			s := NewStatic(http.Dir("."))
			s.Prefix = "/public"
			h.Before(s)

			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(resp.Header().Get("Expires"), convey.ShouldBeBlank)
			convey.So(resp.Body.Len() > 0, convey.ShouldBeTrue)
		})
		convey.Convey("test static redirect", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://127.0.0.1:4000", nil)
			convey.So(err, convey.ShouldBeNil)
			h := New()
			s := NewStatic(http.Dir("."))
			h.Before(s)

			h.ServeHTTP(resp, req)

			convey.So(resp.Code, convey.ShouldEqual, http.StatusFound)
		})
	})
}
