/*
 * Copyright 2015 Xuyuan Pang
 * Author: Pang Xuyuan <xuyuanp # gmail dot com>
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

import "net/http"

type Hador struct {
	Router
	*FilterChain
	Logger Logger
}

func New() *Hador {
	router := NewRouter()
	return &Hador{
		Router:      router,
		FilterChain: NewFilterChain(router),
		Logger:      defaultLogger,
	}
}

func Default() *Hador {
	h := New()
	h.Before(NewLogFilter(h.Logger))
	h.Before(NewRecoveryFilter(h.Logger))
	return h
}

func (h *Hador) Run(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *Hador) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rw := NewResponseWriter(w)
	ctx := &Context{
		Request:  req,
		Response: rw,
		Params:   make(Params),
	}
	h.FilterChain.ServeHTTP(ctx)
}
