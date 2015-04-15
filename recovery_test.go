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

func TestRecovery(t *testing.T) {
	convey.Convey("Given a new Hodor h", t, func() {
		h := New()
		buf := make([]byte, 4096)
		writer := bytes.NewBuffer(buf)
		logger := &logger{
			Logger: log.New(writer, "", 0),
		}
		h.Before(NewRecoveryFilter(logger))
		h.Get("/panic", HandlerFunc(func(ctx *Context) {
			panic("some error")
		}))
		convey.Convey("The panic should be recovery", func() {
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/panic", nil)
			defer func() {
				convey.Convey("Should recover nothing", func() {
					convey.So(recover(), convey.ShouldBeNil)
				})
			}()
			h.ServeHTTP(resp, req)
			convey.So(resp.Code, convey.ShouldEqual, 500)
		})
	})
}
