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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-hodor/hador/swagger"
)

// Hador struct
type Hador struct {
	Router
	*FilterChain
	Logger   Logger
	root     *Node
	Document *swagger.Document
}

// New creates new Hador instance
func New() *Hador {
	h := &Hador{
		Logger: defaultLogger,
		root:   NewNode("", 0),
		Document: &swagger.Document{
			Swagger:     "2.0.0",
			Definitions: swagger.Definitions{},
			Tags:        []swagger.Tag{},
			Responses:   swagger.Responses{},
			Parameters:  map[string]swagger.Parameter{},
			Consumes:    []string{},
			Produces:    []string{},
		},
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

func (h *Hador) travel() swagger.Paths {
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

// Swagger setups swagger config
func (h *Hador) Swagger(config swagger.Config) {
	h.Document.Paths = h.travel()
	h.Get(config.SwaggerPath, HandlerFunc(func(ctx *Context) {
		ctx.Response.Header().Set("Content-Type", "application/json; charset-utf8")
		json.NewEncoder(ctx.Response).Encode(h.Document)
	}))

	s := NewStatic(http.Dir(config.SwaggerUIFilePath))
	s.Prefix = config.SwaggerUIPrefix
	h.Before(s)
}

func (h *Hador) DocHost(host string) *Hador {
	h.Document.Host = host
	return h
}

func (h *Hador) DocBasePath(path string) *Hador {
	h.Document.BasePath = path
	return h
}

func (h *Hador) DocDefinition(model interface{}) *Hador {
	h.Document.Definitions.AddModelFrom(model)
	return h
}

func (h *Hador) DocInfo(title, description, version, termsOfServeice string) *Hador {
	h.Document.Info.Title = title
	h.Document.Info.Description = description
	h.Document.Info.Version = version
	h.Document.Info.TermsOfService = termsOfServeice
	return h
}

func (h *Hador) DocInfoContace(name, url, email string) *Hador {
	h.Document.Info.Contact = &swagger.Contact{
		Name:  name,
		URL:   url,
		Email: email,
	}
	return h
}

func (h *Hador) DocInfoLicense(name, url string) *Hador {
	h.Document.Info.License = &swagger.License{
		Name: name,
		URL:  url,
	}
	return h
}

func (h *Hador) DocConsumes(mimeTypes ...string) *Hador {
	h.Document.Consumes = mimeTypes
	return h
}

func (h *Hador) DocProduces(mimeTypes ...string) *Hador {
	h.Document.Produces = mimeTypes
	return h
}

func (h *Hador) DocTag(name, description string) *Hador {
	h.Document.Tags = append(h.Document.Tags,
		swagger.Tag{Name: name, Description: description})
	return h
}
