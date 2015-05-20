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
	h       *Hador
	parent  *Node
	path    string
	handler Handler
	method  string

	DocIgnored bool
	operation  swagger.Operation
}

// NewLeaf creates new Leaf instance
func NewLeaf(parent *Node, method string, handler Handler) *Leaf {
	l := &Leaf{
		h:       parent.h,
		parent:  parent,
		path:    parent.Path(),
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

// DocIgnore sets if this route would be ignored in document
func (l *Leaf) DocIgnore(ignored bool) *Leaf {
	l.DocIgnored = ignored
	return l
}

// DocTags sets tags of this route
func (l *Leaf) DocTags(tags ...string) *Leaf {
	l.operation.Tags = tags
	return l
}

// DocSummary sets summary of this route
func (l *Leaf) DocSummary(sum string) *Leaf {
	l.operation.Summary = sum
	return l
}

// DocDescription sets description of this route
func (l *Leaf) DocDescription(desc string) *Leaf {
	l.operation.Description = desc
	return l
}

// DocOperationID sets operation id of this route
func (l *Leaf) DocOperationID(oid string) *Leaf {
	l.operation.OperationID = oid
	return l
}

// DocProduces sets produce mime types of this route
func (l *Leaf) DocProduces(mimeTypes ...string) *Leaf {
	l.operation.Produces = mimeTypes
	return l
}

// DocConsumes sets consume mime types of this route
func (l *Leaf) DocConsumes(mimeTypes ...string) *Leaf {
	l.operation.Consumes = mimeTypes
	return l
}

// DocDeprecated sets if this route is deprecated
func (l *Leaf) DocDeprecated(d bool) *Leaf {
	l.operation.Deprecated = d
	return l
}

// DocSchemes set schemes of this route
func (l *Leaf) DocSchemes(schemes ...string) *Leaf {
	l.operation.Schemes = schemes
	return l
}

// DocResponse sets response of this route
func (l *Leaf) DocResponse(code string, resp swagger.Response) *Leaf {
	if l.operation.Responses == nil {
		l.operation.Responses = make(swagger.Responses)
	}
	l.operation.Responses[code] = resp
	return l
}

// DocResponseSimple sets simple response which only contains code and description info of this route
func (l *Leaf) DocResponseSimple(code, desc string) *Leaf {
	resp := swagger.Response{Description: desc}
	l.DocResponse(code, resp)
	return l
}

// DocResponseRef sets response model ref of this route
func (l *Leaf) DocResponseRef(code, desc, ref string) *Leaf {
	resp := swagger.Response{
		Description: desc,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: ref}},
	}
	l.DocResponse(code, resp)
	return l
}

// DocResponseModel sets response model of this route
func (l *Leaf) DocResponseModel(code, desc string, model interface{}) *Leaf {
	l.h.DocDefinition(model)
	resp := swagger.Response{
		Description: desc,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: "#/definitions/" + reflect.TypeOf(model).String()}},
	}
	l.DocResponse(code, resp)
	return l
}

// DocParameter sets parameter of this route
func (l *Leaf) DocParameter(param swagger.Parameter) *Leaf {
	if l.operation.Parameters == nil {
		l.operation.Parameters = make([]swagger.Parameter, 0)
	}
	l.operation.Parameters = append(l.operation.Parameters, param)
	return l
}

// DocPathParameter sets path parameter of this route
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

// DocQueryParameter sets query parameter of this route
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

// DocMultiQueryParameter sets multi query parameter of this route
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

// DocBodyRefParameter set body ref parameter of this route
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

// DocBodyParameter set body model parameter of this route
func (l *Leaf) DocBodyParameter(paramName, desc string, model interface{}, required bool) *Leaf {
	l.h.DocDefinition(model)
	param := swagger.Parameter{
		Name:        paramName,
		In:          "body",
		Description: desc,
		Required:    required,
		Schema:      &swagger.Schema{Reference: swagger.Reference{Ref: "#/definitions/" + reflect.TypeOf(model).String()}},
	}
	l.DocParameter(param)
	return l
}

// DocSecurity set security of this route
func (l *Leaf) DocSecurity(name string, scopes []string) *Leaf {
	if l.operation.Security == nil {
		l.operation.Security = make(swagger.Security)
	}
	l.operation.Security[name] = scopes
	return l
}
