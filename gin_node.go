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

type nodeType int

const (
	static nodeType = iota
	param
	matchAll
)

type node struct {
	segment  string
	indices  string
	children []*node
	ntype    nodeType
	leaves   map[Method]*Leaf
}

func (n *node) AddRoute(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf {
	if len(pattern) == 0 || pattern[0] != '/' {
		panic("pattern should start with '/', pattern: " + pattern)
	}
	if handler == nil {
		panic("handler should NOT be nil")
	}
	return n.addRoute(method, pattern, parseHandler(handler), filters...)
}

func min(first, second int) int {
	if first < second {
		return first
	}
	return second
}

func (n *node) addRoute(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	if len(n.segment) == 0 {
		return n.insertChild(method, pattern, handler, filters...)
	}

	// find longest matched prefix
	max := min(len(n.segment), len(pattern))
	i := 0
	for i < max && pattern[i] == n.segment[i] {
		i++
	}
	// if shorter than n.segment, split current node and insert the reset pattern to new current node.
	if i < max {
		n.splitAt(i)
		return n.insertChild(method, pattern[i:], handler, filters...)
	}

	// if pattern is a prefix of n.segment, split current node, current node handle this route.
	if i == len(pattern) {
		n.splitAt(i)
		return n.handle(method, handler, filters...)
	}

	// if n.segment is a prefix of pattern, insert the reset pattern to current node
	return n.insertChild(method, pattern[i:], handler, filters...)
}

func (n *node) splitAt(index int) {
	next := &node{
		segment:  n.segment[index:],
		indices:  n.indices,
		children: n.children,
		leaves:   n.leaves,
		ntype:    n.ntype,
	}
	n.indices = n.segment[index : index+1]
	n.segment = n.segment[:index]
	n.children = []*node{next}
	n.leaves = nil
}

func (n *node) insertChild(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	if len(n.segment) != 0 {
		n.indices += pattern[:1]
		child := &node{}
		leaf := child.insertChild(method, pattern, handler, filters...)
		n.children = append(n.children, child)
		return leaf
	}
	max := len(pattern)
	i := 0
	for i < max && pattern[i] != '{' {
		i++
	}

	if i == max {
		n.segment = pattern
		n.ntype = static
		n.indices = ""
		n.children = nil
		return n.handle(method, handler, filters...)
	}

	return nil
}

func (n *node) handle(method Method, handler Handler, filters ...Filter) *Leaf {
	if _, ok := n.leaves[method]; ok {
		panic("route has been registered")
	}
	l := NewLeaf(&Node{}, method, handler)
	if n.leaves == nil {
		n.leaves = make(map[Method]*Leaf)
	}
	n.leaves[method] = l
	l.AddFilters(filters...)
	return l
}

func (n *node) find(method Method, path string) *Leaf {
	switch n.ntype {
	case static:
		return n.findStatic(method, path)
	}
	return nil
}

func (n *node) findStatic(method Method, path string) *Leaf {
	if len(path) < len(n.segment) {
		return nil
	}
	if path == n.segment {
		if n.leaves != nil {
			return n.leaves[method]
		}
	}
	if path[:len(n.segment)] != n.segment {
		return nil
	}
	c := path[len(n.segment)]
	for i, ind := range n.indices {
		if ind == rune(c) {
			return n.children[i].find(method, path[len(n.segment):])
		}
	}
	return nil
}
