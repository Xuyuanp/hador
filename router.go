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

// Router interface
type Router interface {
	Handler

	Options(string, Handler) Beforer
	Get(string, Handler) Beforer
	Head(string, Handler) Beforer
	Post(string, Handler) Beforer
	Put(string, Handler) Beforer
	Delete(string, Handler) Beforer
	Trace(string, Handler) Beforer
	Connect(string, Handler) Beforer
	Patch(string, Handler) Beforer
	Any(string, Handler) Beforer

	AddRoute(string, string, Handler) Beforer

	Group(string, func(Router)) Beforer
}

// NewRouter creates Router interface
func NewRouter(h *Hador) Router {
	return newNode(h, "", 0)
}
