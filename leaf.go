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

// Leaf struct
type Leaf struct {
	*FilterChain
	parent  *Node
	path    string
	handler Handler
	method  string
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
