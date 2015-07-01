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

// Document struct
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

// DocHost sets dochost of document
func (doc *Document) DocHost(host string) *Document {
	doc.Host = host
	return doc
}

// DocBasePath sets basepath of document
func (doc *Document) DocBasePath(path string) *Document {
	doc.BasePath = path
	return doc
}

// DocDefinition adds model definition
func (doc *Document) DocDefinition(model interface{}) *Document {
	doc.Definitions.AddModelFrom(model)
	return doc
}

// DocInfo sets info of document
func (doc *Document) DocInfo(title, description, version, termsOfServeice string) *Document {
	doc.Info.Title = title
	doc.Info.Description = description
	doc.Info.Version = version
	doc.Info.TermsOfService = termsOfServeice
	return doc
}

// DocInfoContace sets info contace of document
func (doc *Document) DocInfoContace(name, url, email string) *Document {
	doc.Info.Contact = &Contact{
		Name:  name,
		URL:   url,
		Email: email,
	}
	return doc
}

// DocInfoLicense sets info license of document
func (doc *Document) DocInfoLicense(name, url string) *Document {
	doc.Info.License = &License{
		Name: name,
		URL:  url,
	}
	return doc
}

// DocConsumes sets consumes of document
func (doc *Document) DocConsumes(mimeTypes ...string) *Document {
	doc.Consumes = mimeTypes
	return doc
}

// DocProduces sets produces of document
func (doc *Document) DocProduces(mimeTypes ...string) *Document {
	doc.Produces = mimeTypes
	return doc
}

// DocTag adds tag to document
func (doc *Document) DocTag(name, description string) *Document {
	doc.Tags = append(doc.Tags,
		Tag{Name: name, Description: description})
	return doc
}
