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
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Xuyuanp/hador/swagger"
)

// Hador struct
type Hador struct {
	Router
	*FilterChain
	Logger Logger
	root   *node

	ctxPool  sync.Pool
	respPool sync.Pool

	Document *swagger.Document
}

// New creates new Hador instance
func New() *Hador {
	h := &Hador{
		Logger: defaultLogger,
		Document: &swagger.Document{
			Swagger:     "2.0.0",
			Definitions: swagger.GlobalDefinitions,
			Tags:        []swagger.Tag{},
			Responses:   swagger.Responses{},
			Parameters:  map[string]swagger.Parameter{},
			Consumes:    []string{},
			Produces:    []string{},
		},
	}
	h.root = &node{}
	h.Router = RouterFunc(h.root.AddRoute)
	h.FilterChain = NewFilterChain(h.root)

	h.ctxPool.New = func() interface{} {
		return newContext(h.Logger)
	}
	h.respPool.New = func() interface{} {
		return NewResponseWriter(nil)
	}

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

// RunTLS starts serving HTTPS request.
func (h *Hador) RunTLS(addr, sertFile, keyFile string) error {
	h.Logger.Info("Listening on %s", addr)
	return http.ListenAndServeTLS(addr, sertFile, keyFile, h)
}

func (h *Hador) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := h.respPool.Get().(*responseWriter)
	resp.reset(w)
	defer h.respPool.Put(resp)

	ctx := h.ctxPool.Get().(*Context)
	ctx.reset(resp, req)
	defer h.ctxPool.Put(ctx)

	h.Serve(ctx)
}

// Serve implements Handler interface
func (h *Hador) Serve(ctx *Context) {
	ctx.SetHeader("Server", "Hador/"+Version)
	h.FilterChain.Serve(ctx)
}

// AddFilters reuses FilterChain's AddFilters method and returns self
func (h *Hador) AddFilters(filters ...Filter) *Hador {
	h.FilterChain.AddFilters(filters...)
	return h
}

func (h *Hador) travel() []*Leaf {
	llist := list.New()
	h.root.travel(llist)

	leaves := make([]*Leaf, llist.Len())
	i := 0
	for e := llist.Front(); e != nil; e = e.Next() {
		leaves[i] = e.Value.(*Leaf)
		i++
	}
	return leaves
}

func (h *Hador) travelPaths() swagger.Paths {
	spaths := make(swagger.Paths)
	leaves := h.travel()
	for _, leaf := range leaves {
		if leaf.DocIgnored || leaf.method == "ANY" {
			continue
		}

		spath, ok := spaths[leaf.Path()]
		if !ok {
			spath = make(swagger.Path)
			spaths[leaf.Path()] = spath
		}

		spath[strings.ToLower(leaf.Method().String())] = *leaf.SwaggerOperation()
	}
	return spaths
}

// SwaggerHandler returns swagger json api handler
func (h *Hador) SwaggerHandler() Handler {
	return HandlerFunc(func(ctx *Context) {
		h.Document.Paths = h.travelPaths()
		ctx.SetHeader("Content-Type", "application/json; charset-utf8")
		ctx.RenderJSON(h.Document)
	})
}

// Swagger setups swagger config, returns json API path Leaf
func (h *Hador) Swagger(config swagger.Config) *Leaf {
	// handle API path
	leaf := h.Get(config.APIPath, h.SwaggerHandler()).DocIgnore(!config.SelfDocEnabled)

	// serve swagger-ui file
	if config.UIFilePath != "" {
		s := NewStatic(http.Dir(config.UIFilePath))
		s.Prefix = config.UIPrefix
		h.Before(s)
	}

	return leaf
}

// SwaggerDocument returns swagger.Document of this Hador.
func (h *Hador) SwaggerDocument() *swagger.Document {
	return h.Document
}

func (h *Hador) showGraph() {
	leaves := h.travel()
	for _, l := range leaves {
		fmt.Println(l.Path())
	}
}

func (h *Hador) _showGraph() {
	h.root._travel("")
}
