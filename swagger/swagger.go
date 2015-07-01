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

// Package swagger provides swagger specification models
package swagger

// Security type
type Security map[string][]string

// Scopes type
type Scopes map[string]string

// Paths type
type Paths map[string]Path

// Path type
type Path map[string]Operation

// Responses type
type Responses map[string]Response

type Parameters []Parameter

// Headers type
type Headers map[string]Header

// Example type
type Example map[string]interface{}

// SecurityDefinitons type
type SecurityDefinitons map[string]SecurityDefiniton

// Contact struct
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License struct
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Info struct
type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Reference struct
type Reference struct {
	Ref string `json:"$ref,omitempty"`
}

// Items struct
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

// Header struct
type Header struct {
	Items
	Description string `json:"description,omitempty"`
}

// XML struct
type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

// Schema struct
type Schema struct {
	Reference
	Type       string           `json:"type,omitempty"`
	Required   []string         `json:"required,omitempty"`
	Properties map[string]Items `json:"properties,omitempty"`
}

// Response struct
type Response struct {
	Description string   `json:"description"`
	Schema      *Schema  `json:"schema,omitempty"`
	Headers     *Headers `json:"headers,omitempty"`
	Example     Example  `json:"example,omitempty"`
}

// Parameter struct
type Parameter struct {
	Items
	Name        string  `json:"name,omitempty"`
	In          string  `json:"in"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

// SecurityDefiniton struct
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

// Tag struct
type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs struct
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}
