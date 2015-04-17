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

import (
	"container/list"
	"net/http"
	"strings"

	"github.com/go-hodor/hador/swagger"
)

// Hador struct
type Hador struct {
	Router
	*FilterChain
	Logger Logger
	root   *Node
}

// New creates new Hador instance
func New() *Hador {
	h := &Hador{
		Logger: defaultLogger,
		root:   NewNode("", 0),
	}
	h.Router = h.root
	h.FilterChain = NewFilterChain(h.Router)
	return h
}

// Default creates Hador instance with default filters(LogFilter, RecoveryFilter)
func Default() *Hador {
	h := New()
	h.Before(NewLogFilter(h.Logger))
	h.Before(NewRecoveryFilter(h.Logger))
	return h
}

// Run starts serving HTTP request
func (h *Hador) Run(addr string) error {
	h.Logger.Info("Listening on %s", addr)
	return http.ListenAndServe(addr, h)
}

func (h *Hador) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(w, req, h.Logger)
	h.Serve(ctx)
}

// Serve implements Handler interface
func (h *Hador) Serve(ctx *Context) {
	h.FilterChain.Serve(ctx)
}

func (h *Hador) Travel() swagger.Paths {
	llist := list.New()
	h.root.travel(llist)

	spaths := make(swagger.Paths)

	e := llist.Front()
	for e != nil {
		leaf := e.Value.(*Leaf)

		spath, ok := spaths[leaf.Path()]
		if !ok {
			spath = make(swagger.Path)
			spaths[leaf.Path()] = spath
		}

		spath[strings.ToLower(leaf.Method())] = leaf.operation

		e = e.Next()
	}
	return spaths
}
