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
	"reflect"

	"github.com/go-hodor/hador/swagger"
)

// Leaf struct
type Leaf struct {
	*FilterChain
	parent  *Node
	path    string
	handler Handler
	method  string

	DocIgnored bool
	operation  swagger.Operation
}

// NewLeaf creates new Leaf instance
func NewLeaf(method string, handler Handler) *Leaf {
	l := &Leaf{
		method:  method,
		handler: handler,
	}
	l.FilterChain = NewFilterChain(l.handler)
	return l
}

// Path returns the full path from root to the parent node
func (l *Leaf) Path() string {
	return l.path
}

// Method returns method of Leaf
func (l *Leaf) Method() string {
	return l.method
}

// Handler returns handler of Leaf
func (l *Leaf) Handler() Handler {
	return l.handler
}

// Parent returns parent node of leaf
func (l *Leaf) Parent() *Node {
	return l.parent
}

// AddFilters add filters into FilterChain
func (l *Leaf) AddFilters(filters ...Filter) *Leaf {
	l.FilterChain.AddFilters(filters...)
	return l
}

func (l *Leaf) DocIgnore(ignored bool) *Leaf {
	l.DocIgnored = ignored
	return l
}

func (l *Leaf) DocTags(tags ...string) *Leaf {
	l.operation.Tags = tags
	return l
}

func (l *Leaf) DocSummary(sum string) *Leaf {
	l.operation.Summary = sum
	return l
}

func (l *Leaf) DocDescription(desc string) *Leaf {
	l.operation.Description = desc
	return l
}

func (l *Leaf) DocOperationID(oid string) *Leaf {
	l.operation.OperationID = oid
	return l
}

func (l *Leaf) DocProduces(mimeTypes ...string) *Leaf {
	l.operation.Produces = mimeTypes
	return l
}

func (l *Leaf) DocConsumes(mimeTypes ...string) *Leaf {
	l.operation.Consumes = mimeTypes
	return l
}

func (l *Leaf) DocDeprecated(d bool) *Leaf {
	l.operation.Deprecated = d
	return l
}

func (l *Leaf) DocSchemes(schemes []string) *Leaf {
	l.operation.Schemes = schemes
	return l
}

func (l *Leaf) DocResponse(code string, resp swagger.Response) *Leaf {
	if l.operation.Responses == nil {
		l.operation.Responses = make(swagger.Responses)
	}
	l.operation.Responses[code] = resp
	return l
}

func (l *Leaf) DocResponseSimple(code, desc string) *Leaf {
	resp := swagger.Response{Description: desc}
	l.DocResponse(code, resp)
	return l
}

func (l *Leaf) DocResponseRef(code, desc, ref string) *Leaf {
	resp := swagger.Response{
		Description: desc,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: ref}},
	}
	l.DocResponse(code, resp)
	return l
}

func (l *Leaf) DocResponseModel(code, desc string, model interface{}) *Leaf {
	resp := swagger.Response{
		Description: desc,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: "#/definitions/" + reflect.TypeOf(model).Name()}},
	}
	l.DocResponse(code, resp)
	return l
}

func (l *Leaf) DocParameter(param swagger.Parameter) *Leaf {
	if l.operation.Parameters == nil {
		l.operation.Parameters = make([]swagger.Parameter, 0)
	}
	l.operation.Parameters = append(l.operation.Parameters, param)
	return l
}

func (l *Leaf) DocPathParameter(paramName, paramType, desc string, required bool) *Leaf {
	param := swagger.Parameter{
		Name:        paramName,
		In:          "path",
		Description: desc,
		Required:    required,
		Items: swagger.Items{
			Type: paramType,
		},
	}
	l.DocParameter(param)
	return l
}

func (l *Leaf) DocQueryParameter(paramName, paramType, desc string, required bool) *Leaf {
	param := swagger.Parameter{
		Name:        paramName,
		In:          "query",
		Description: desc,
		Required:    required,
		Items: swagger.Items{
			Type: paramType,
		},
	}
	l.DocParameter(param)
	return l
}

func (l *Leaf) DocMultiQueryParameter(paramName, paramType, desc string, required bool) *Leaf {
	param := swagger.Parameter{
		Name:        paramName,
		In:          "query",
		Description: desc,
		Required:    required,
		Items: swagger.Items{
			Type: "array",
			Items: &swagger.Items{
				Type: paramType,
			},
			CollectionFormat: "multi",
		},
	}
	l.DocParameter(param)
	return l
}

func (l *Leaf) DocBodyRefParameter(paramName, desc string, ref string, required bool) *Leaf {
	param := swagger.Parameter{
		Name:        paramName,
		In:          "body",
		Description: desc,
		Required:    required,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: ref}},
	}
	l.DocParameter(param)
	return l
}

func (l *Leaf) DocBodyParameter(paramName, desc string, model interface{}, required bool) *Leaf {
	param := swagger.Parameter{
		Name:        paramName,
		In:          "body",
		Description: desc,
		Required:    required,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: "#/definitions/" + reflect.TypeOf(model).Name()}},
	}
	l.DocParameter(param)
	return l
}

func (l *Leaf) DocSecurity(name string, scopes []string) *Leaf {
	if l.operation.Security == nil {
		l.operation.Security = make(swagger.Security)
	}
	l.operation.Security[name] = scopes
	return l
}
