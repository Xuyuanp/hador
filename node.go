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

import (
	"container/list"
	"fmt"
	"net/http"
	"regexp"
)

type nodeType int

const (
	static nodeType = iota
	param
	matchAll
)

type node struct {
	parent     *node
	segment    string
	indices    string
	children   []*node
	paramChild *node
	ntype      nodeType
	leaves     map[Method]*Leaf

	paramName     string
	paramReg      *regexp.Regexp
	paramDataType string
	paramDesc     string
}

func (n *node) AddRoute(method Method, pattern string, handler interface{}, filters ...Filter) *Leaf {
	if len(pattern) == 0 || pattern[0] != '/' {
		panic("pattern should start with '/', pattern: " + pattern)
	}
	if len(pattern) > 1 && pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}
	if handler == nil {
		panic("handler should NOT be nil")
	}
	for _, m := range Methods {
		if m == method {
			return n.addRoute(method, pattern, parseHandler(handler), filters...)
		}
	}
	panic("unknown method: " + method)
}

func min(first, second int) int {
	if first < second {
		return first
	}
	return second
}

func (n *node) addRoute(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	if len(n.segment) == 0 {
		return n.init(method, pattern, handler, filters...)
	}

	if n.ntype == static {
		// find longest matched prefix
		max := min(len(n.segment), len(pattern))
		i := 0
		for i < max && pattern[i] == n.segment[i] {
			i++
		}
		n.splitAt(i)
		return n.insertChild(method, pattern[i:], handler, filters...)
	}

	if n.ntype == param {
		i, max := 0, len(pattern)
		for i < max && pattern[i] != '}' {
			i++
		}
		if i == max {
			panic("missing '}'")
		}
		if n.segment != pattern[:i+1] {
			panic("conflict param node")
		}
		if i < max-1 && pattern[i+1] != '/' {
			panic("'}' should be before '/'")
		}
		return n.insertChild(method, pattern[i+1:], handler, filters...)
	}
	return nil
}

func (n *node) splitAt(index int) {
	if index >= len(n.segment) {
		return
	}
	next := &node{
		parent:     n,
		segment:    n.segment[index:],
		indices:    n.indices,
		children:   n.children,
		leaves:     n.leaves,
		ntype:      n.ntype,
		paramChild: n.paramChild,
	}
	if next.children != nil {
		for _, ch := range next.children {
			ch.parent = next
		}
	}
	if next.paramChild != nil {
		next.paramChild.parent = next
	}
	n.indices = n.segment[index : index+1]
	n.segment = n.segment[:index]
	n.children = []*node{next}
	n.paramChild = nil
	n.leaves = nil
}

func (n *node) insertChild(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	if len(pattern) == 0 {
		return n.handle(method, handler, filters...)
	}
	if pattern[0] == '{' {
		return n.insertParamChild(method, pattern, handler, filters...)
	}
	return n.insertStaticChild(method, pattern, handler, filters...)
}

func (n *node) insertStaticChild(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	for i, ind := range n.indices {
		if ind == rune(pattern[0]) {
			return n.children[i].addRoute(method, pattern, handler, filters...)
		}
	}
	n.indices += pattern[:1]
	child := &node{parent: n}
	n.children = append(n.children, child)
	return child.addRoute(method, pattern, handler, filters...)
}

func (n *node) insertParamChild(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	if n.paramChild == nil {
		n.paramChild = &node{parent: n}
	}
	return n.paramChild.addRoute(method, pattern, handler, filters...)
}

func (n *node) init(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	if pattern[0] == '{' {
		return n.initParam(method, pattern, handler, filters...)
	}
	return n.initStatic(method, pattern, handler, filters...)
}

func (n *node) initStatic(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	i, max := 0, len(pattern)
	for i < max && pattern[i] != '{' {
		i++
	}
	if i < max && i > 0 && pattern[i-1] != '/' {
		panic("'{' should be after '/'")
	}

	n.segment = pattern[:i]
	n.ntype = static
	n.indices = ""
	n.children = nil
	n.leaves = nil
	return n.insertChild(method, pattern[i:], handler, filters...)
}

func (n *node) initParam(method Method, pattern string, handler Handler, filters ...Filter) *Leaf {
	name, regstr, dataType, desc, rest := readParam(pattern)
	if len(name) == 0 {
		panic("empty param name")
	}
	n.paramName = name
	if len(regstr) > 0 {
		n.paramReg = regexp.MustCompile(regstr)
	}
	n.paramDataType = dataType
	n.paramDesc = desc

	n.ntype = param
	n.segment = pattern[:len(pattern)-len(rest)]
	n.indices = ""
	n.children = nil
	n.leaves = nil
	return n.insertChild(method, rest, handler, filters...)
}

func (n *node) handle(method Method, handler Handler, filters ...Filter) *Leaf {
	if _, ok := n.leaves[method]; ok {
		panic("route has been registered")
	}
	l := NewLeaf(n, method, handler)
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
	case param:
		return n.findParam(method, path)
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
		return nil
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
	if n.paramChild != nil {
		return n.paramChild.find(method, path[len(n.segment):])
	}
	return nil
}

func (n *node) findParam(method Method, path string) *Leaf {
	i, max := 0, len(path)
	for i < max && path[i] != '/' {
		i++
	}
	if i == max {
		if n.leaves != nil {
			return n.leaves[method]
		}
		return nil
	}
	c := path[i]
	for index, ind := range n.indices {
		if ind == rune(c) {
			return n.children[index].find(method, path[i:])
		}
	}
	if n.paramChild != nil {
		return n.paramChild.find(method, path[i:])
	}
	return nil
}

func (n *node) Serve(ctx *Context) {
	// ctx.Logger.Debug("%s", ctx.Request.RequestURI)
	switch n.ntype {
	case static:
		n.serveStatic(ctx)
	case param:
		n.serveParam(ctx)
	}
}

func (n *node) serveParam(ctx *Context) {
	path := ctx.path
	i, max := 0, len(path)
	for i < max && path[i] != '/' {
		i++
	}
	if i == max {
		ctx.Params()[n.segment[1:len(n.segment)-1]] = path[:i]
		n.doServe(ctx)
		return
	}
	c := path[i]
	for index, ind := range n.indices {
		if ind == rune(c) {
			ctx.path = ctx.path[i:]
			ctx.Params()[n.segment[1:len(n.segment)-1]] = path[:i]
			n.children[index].Serve(ctx)
			return
		}
	}
	if n.paramChild != nil {
		ctx.path = ctx.path[1:]
		ctx.Params()[n.segment[1:len(n.segment)-1]] = path[:i]
		n.paramChild.Serve(ctx)
		return
	}
	return
}

func (n *node) serveStatic(ctx *Context) {
	path := ctx.path
	if len(path) < len(n.segment) {
		ctx.OnError(http.StatusNotFound)
		return
	}
	i, seglen := 0, len(n.segment)
	for i < seglen && n.segment[i] == path[i] {
		i++
	}
	if i < seglen {
		ctx.OnError(http.StatusNotFound)
		return
	}
	if i == len(path) {
		n.doServe(ctx)
		return
	}
	c := path[seglen]
	for index, ind := range n.indices {
		if ind == rune(c) {
			ctx.path = ctx.path[seglen:]
			n.children[index].Serve(ctx)
			return
		}
	}
	if n.paramChild != nil {
		ctx.path = ctx.path[len(n.segment):]
		n.paramChild.Serve(ctx)
		return
	}
	ctx.OnError(http.StatusNotFound)
}

func (n *node) doServe(ctx *Context) {
	// 404 not found
	if len(n.leaves) == 0 {
		ctx.OnError(http.StatusNotFound)
		return
	}
	// method matches
	if l, ok := n.leaves[Method(ctx.Request.Method)]; ok {
		l.Serve(ctx)
		return
	}
	// 405 method not allowed
	methods := make([]Method, len(n.leaves))
	i := 0
	for m := range n.leaves {
		methods[i] = m
		i++
	}
	ctx.OnError(http.StatusMethodNotAllowed, methods)
}

func (n *node) travel(llist *list.List) {
	for _, l := range n.leaves {
		llist.PushBack(l)
	}

	for _, child := range n.children {
		child.travel(llist)
	}
	if n.paramChild != nil {
		n.paramChild.travel(llist)
	}
}

func (n *node) path() string {
	var path string
	if n.ntype == static {
		path = n.segment
	} else if n.ntype == param {
		path = "{" + n.paramName + "}"
	}
	if n.parent != nil {
		return n.parent.path() + path
	}
	return path
}

func (n *node) _travel(path string) {
	path += n.segment
	for m, l := range n.leaves {
		if path != l.path {
			fmt.Printf("%s %s ===== %s\n", m, path, l.Path())
		}
	}
	for _, child := range n.children {
		child._travel(path)
	}
	if n.paramChild != nil {
		n.paramChild._travel(path)
	}
}

func readParam(pattern string) (name, regstr, dataType, desc, rest string) {
	dataType = "string"
	field, rest, end := readField(pattern[1:])
	name = field
	if end || len(rest) == 0 {
		return
	}

	field, rest, end = readField(rest)
	regstr = field
	if end || len(rest) == 0 {
		return
	}

	field, rest, end = readField(rest)
	if len(field) > 0 {
		dataType = field
	}
	if end || len(rest) == 0 {
		return
	}

	field, rest, end = readField(rest)
	desc = field
	if !end {
		panic("missing '}'")
	}

	return
}

func readField(pattern string) (field, rest string, end bool) {
	i, max := 0, len(pattern)
	for i < max {
		if pattern[i] == ':' {
			return pattern[:i], pattern[i+1:], false
		} else if pattern[i] == '}' {
			if i < max-1 && pattern[i+1] != '/' {
				panic("'}' should be before '/' or at the end")
			}
			return pattern[:i], pattern[i+1:], true
		}
		i++
	}
	panic("'}' should be before '/' or at the end")
}
