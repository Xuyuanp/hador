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

	Options(pattern string, h interface{}, filters ...Filter) *Leaf
	Get(pattern string, h interface{}, filters ...Filter) *Leaf
	Head(pattern string, h interface{}, filters ...Filter) *Leaf
	Post(pattern string, h interface{}, filters ...Filter) *Leaf
	Put(pattern string, h interface{}, filters ...Filter) *Leaf
	Delete(pattern string, h interface{}, filters ...Filter) *Leaf
	Trace(pattern string, h interface{}, filters ...Filter) *Leaf
	Connect(pattern string, h interface{}, filters ...Filter) *Leaf
	Patch(pattern string, h interface{}, filters ...Filter) *Leaf
	Any(pattern string, h interface{}, filters ...Filter) *Leaf

	AddController(pattern string, controller ControllerInterface, filters ...Filter)

	AddRoute(pattern string, method string, h interface{}, filters ...Filter) *Leaf

	Group(pattern string, fn func(Router), filters ...Filter)
}
