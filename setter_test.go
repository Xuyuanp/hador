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
)

func TestSetter(t *testing.T) {
	cases := []struct {
		method  Method
		pattern string
		handler interface{}

		path string
		code int
		body string
	}{
		{GET, "/hello", emptyHandler, "/hello", 200, ""},
		{GET, "/hello/{name}", newParamHandler("name", ""), "/hello/jack", 200, "jack"},
	}

	h := New()

	for _, c := range cases {
		h.Route().Method(c.method).Pattern(c.pattern).Handler(c.handler)
	}

	for _, c := range cases {
		req, _ := http.NewRequest(c.method.String(), c.path, nil)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		if resp.Code != c.code || resp.Body.String() != c.body {
			t.Errorf("%s %s Failed: %d => %d %s => %s", c.code, resp.Code, c.body, resp.Body.String())
		}
	}
}
