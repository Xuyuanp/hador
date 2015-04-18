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

type Security map[string][]string

type Scopes map[string]string

type Paths map[string]Path

type Path map[string]Operation

type Responses map[string]Response

type Headers map[string]Header

type Example map[string]interface{}

type SecurityDefinitons map[string]SecurityDefiniton

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Reference struct {
	Ref string `json:"$ref,omitempty"`
}

type Items struct {
	Reference
	Items            *Items        `json:"items,omitempty"`
	Type             string        `json:"type,omitempty"`
	Format           string        `json:"format,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Maximum          int           `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          int           `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        int           `json:"maxLength,omitempty"`
	MinLength        int           `json:"minLength,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	MultipleOf       int           `json:"multipleOf,omitempty"`
}

type Header struct {
	Items
	Description string `json:"description,omitempty"`
}

type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

type Schema struct {
	Reference
	Type       string           `json:"type,omitempty"`
	Required   []string         `json:"required,omitempty"`
	Properties map[string]Items `json:"properties,omitempty"`
}

type Response struct {
	Description string   `json:"description"`
	Schema      *Schema  `json:"schema,omitempty"`
	Headers     *Headers `json:"headers,omitempty"`
	Example     Example  `json:"example,omitempty"`
}

type Parameter struct {
	Items
	Name        string  `json:"name,omitempty"`
	In          string  `json:"in"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type Operation struct {
	Reference
	Tags         []string      `json:"tags,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	OperationID  string        `json:"operationId,omitempty"`
	Parameters   []Parameter   `json:"parameters,omitempty"`
	Consumes     []string      `json:"comsumes,omitempty"`
	Produces     []string      `json:"produces,omitempty"`
	Responses    Responses     `json:"responses,omitempty"`
	Schemes      []string      `json:"schemes,omitempty"`
	Deprecated   bool          `json:"deprecated,omitempty"`
	Security     Security      `json:"security,omitempty"`
}

type SecurityDefiniton struct {
	Type             string `json:"type"`
	Description      string `json:"description,omitempty"`
	Name             string `json:"name"`
	In               string `json:"in"`
	Flow             string `json:"flow"`
	AuthorizationURL string `json:"authorizationUrl"`
	TokenURL         string `json:"tokenUrl"`
	Scopes           Scopes `json:"scopes"`
}

type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

type Definitions map[string]Schema

type Document struct {
	Swagger            string               `json:"swagger"`
	Info               Info                 `json:"info"`
	Host               string               `json:"host,omitempty"`
	BasePath           string               `json:"basePath,omitempty"`
	Schemes            []string             `json:"schemes,omitempty"`
	Consumes           []string             `json:"consumes,omitempty"`
	Produces           []string             `json:"produces,omitempty"`
	Paths              Paths                `json:"paths"`
	Definitions        Definitions          `json:"definitions,omitempty"`
	Parameters         map[string]Parameter `json:"parameters,omitempty"`
	Responses          Responses            `json:"responses,omitempty"`
	SecurityDefinitons SecurityDefinitons   `json:"securityDefinitions,omitempty"`
	Security           Security             `json:"security,omitempty"`
	Tags               []Tag                `json:"tags,omitempty"`
	ExternalDocs       *ExternalDocs        `json:"externalDocs,omitempty"`
}
