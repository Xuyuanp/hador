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

package swagger

import "reflect"

// Operation struct
type Operation struct {
	Reference
	Tags         []string      `json:"tags,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	OperationID  string        `json:"operationId,omitempty"`
	Parameters   Parameters    `json:"parameters,omitempty"`
	Consumes     []string      `json:"comsumes,omitempty"`
	Produces     []string      `json:"produces,omitempty"`
	Responses    Responses     `json:"responses,omitempty"`
	Schemes      []string      `json:"schemes,omitempty"`
	Deprecated   bool          `json:"deprecated,omitempty"`
	Security     Security      `json:"security,omitempty"`
}

// DocSumDesc sets summary and description of this operation
func (o *Operation) DocSumDesc(summary, description string) *Operation {
	o.Summary = summary
	o.Description = description
	return o
}

// DocTags sets tags of this Operation
func (o *Operation) DocTags(tags ...string) *Operation {
	o.Tags = tags
	return o
}

// DocProduces sets produce mime types of this Operation
func (o *Operation) DocProduces(mimeTypes ...string) *Operation {
	o.Produces = mimeTypes
	return o
}

// DocConsumes sets consume mime types of this Operation
func (o *Operation) DocConsumes(mimeTypes ...string) *Operation {
	o.Consumes = mimeTypes
	return o
}

// DocDeprecated sets if this Operation is deprecated
func (o *Operation) DocDeprecated(d bool) *Operation {
	o.Deprecated = d
	return o
}

// DocSchemes set schemes of this Operation
func (o *Operation) DocSchemes(schemes ...string) *Operation {
	o.Schemes = schemes
	return o
}

// DocResponse sets response of this Operation
func (o *Operation) DocResponse(code string, resp Response) *Operation {
	if o.Responses == nil {
		o.Responses = make(Responses)
	}
	o.Responses[code] = resp
	return o
}

// DocResponseSimple sets simple response which only contains code and description info of this Operation
func (o *Operation) DocResponseSimple(code string, desc string) *Operation {
	resp := Response{Description: desc}
	o.DocResponse(code, resp)
	return o
}

// DocResponseRef sets response model ref of this Operation
func (o *Operation) DocResponseRef(code, desc, ref string) *Operation {
	resp := Response{
		Description: desc,
		Schema:      &Schema{Reference: Reference{Ref: ref}},
	}
	o.DocResponse(code, resp)
	return o
}

// DocResponseModel sets response model of this Operation
func (o *Operation) DocResponseModel(code string, desc string, model interface{}) *Operation {
	GlobalDefinitions.AddModelFrom(model)
	resp := Response{
		Description: desc,
		Schema: &Schema{
			Reference: Reference{Ref: "#/definitions/" + reflect.TypeOf(model).String()},
		},
	}
	o.DocResponse(code, resp)
	return o
}

// DocParameter sets parameter of this Operation
func (o *Operation) DocParameter(param Parameter) *Operation {
	if o.Parameters == nil {
		o.Parameters = make(Parameters, 0)
	}
	o.Parameters = append(o.Parameters, param)
	return o
}

// DocParameterPath sets path parameter of this Operation
func (o *Operation) DocParameterPath(paramName, paramType, desc string, required bool) *Operation {
	param := Parameter{
		Name:        paramName,
		In:          "path",
		Description: desc,
		Required:    required,
		Items: Items{
			Type: paramType,
		},
	}
	o.DocParameter(param)
	return o
}

// DocParameterQuery sets query parameter of this Operation
func (o *Operation) DocParameterQuery(paramName, paramType, desc string, required bool) *Operation {
	param := Parameter{
		Name:        paramName,
		In:          "query",
		Description: desc,
		Required:    required,
		Items: Items{
			Type: paramType,
		},
	}
	o.DocParameter(param)
	return o
}

// DocParameterMultiQuery sets multi query parameter of this Operation
func (o *Operation) DocParameterMultiQuery(paramName, paramType, desc string, required bool) *Operation {
	param := Parameter{
		Name:        paramName,
		In:          "query",
		Description: desc,
		Required:    required,
		Items: Items{
			Type: "array",
			Items: &Items{
				Type: paramType,
			},
			CollectionFormat: "multi",
		},
	}
	o.DocParameter(param)
	return o
}

// DocParameterBodyRef set body ref parameter of this Operation
func (o *Operation) DocParameterBodyRef(paramName, desc string, ref string, required bool) *Operation {
	param := Parameter{
		Name:        paramName,
		In:          "body",
		Description: desc,
		Required:    required,
		Schema:      &Schema{Reference: Reference{Ref: ref}},
	}
	o.DocParameter(param)
	return o
}

// DocParameterBody set body model parameter of this Operation
func (o *Operation) DocParameterBody(paramName, desc string, model interface{}, required bool) *Operation {
	GlobalDefinitions.AddModelFrom(model)
	param := Parameter{
		Name:        paramName,
		In:          "body",
		Description: desc,
		Required:    required,
		Schema:      &Schema{Reference: Reference{Ref: "#/definitions/" + reflect.TypeOf(model).String()}},
	}
	o.DocParameter(param)
	return o
}

// DocSecurity set security of this Operation
func (o *Operation) DocSecurity(name string, scopes []string) *Operation {
	if o.Security == nil {
		o.Security = make(Security)
	}
	o.Security[name] = scopes
	return o
}
