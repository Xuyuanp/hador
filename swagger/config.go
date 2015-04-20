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

// Config struct
type Config struct {
	// UIFilePath is the location of folder containing swagger-ui index.html file. e.g. swagger-ui/dist
	UIFilePath string

	// UIPrefx is the path where swagger-ui whill be served. e.g. /apidocs
	UIPrefix string

	// APIPath is the path where JSON API is available. e.g. /apidocs.json
	APIPath string

	// CORSDisabled disable CORS filter. False on default.
	CORSDisabled bool

	// SelfDocEnabled enable the swagger-ui path API in doc. False on default.
	SelfDocEnabled bool
}
