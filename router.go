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

	Options(string, Handler, ...Filter) *Leaf
	Get(string, Handler, ...Filter) *Leaf
	Head(string, Handler, ...Filter) *Leaf
	Post(string, Handler, ...Filter) *Leaf
	Put(string, Handler, ...Filter) *Leaf
	Delete(string, Handler, ...Filter) *Leaf
	Trace(string, Handler, ...Filter) *Leaf
	Connect(string, Handler, ...Filter) *Leaf
	Patch(string, Handler, ...Filter) *Leaf
	Any(string, Handler, ...Filter) *Leaf

	AddRoute(string, string, Handler, ...Filter) *Leaf

	Group(string, func(Router), ...Filter)
}
