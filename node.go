/*
 * Copyright 2015 Xuyuan Pang <xuyuanp # gmail dot com>
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
	"strings"
)

func splitSegments(path string) []string {
	path = trimPath(path)
	if len(path) == 0 {
		return []string{}
	}
	return strings.Split(path, "/")
}

func isReg(segment string) bool {
	if regSegmentRegexp.MatchString(segment) {
		return true
	}
	if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
		return true
	}
	return false
}

var regSegmentRegexp = regexp.MustCompile(`\(\?P<.+>.+\)`)

// Node struct
type Node struct {
	h           *Hador
	parent      *Node
	depth       int
	segment     string
	paramName   string
	paramReg    *regexp.Regexp
	rawChildren map[string]*Node
	regChildren []*Node
	leaves      map[Method]*Leaf
}

// NewNode creates new Node instance.
func NewNode(h *Hador, segment string, depth int) *Node {
	paramName, paramReg := resolveSegment(segment)
	n := &Node{
		h:           h,
		segment:     segment,
		depth:       depth,
		paramName:   paramName,
		paramReg:    paramReg,
		rawChildren: make(map[string]*Node),
		regChildren: make([]*Node, 0),
		leaves:      make(map[Method]*Leaf),
	}
	return n
}

func resolveSegment(segment string) (paramName string, paramReg *regexp.Regexp) {
	var splits []string
	if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
		seg := segment[1 : len(segment)-1]
		splits = strings.SplitN(seg, ":", 2)
		if len(splits) == 1 {
			splits = append(splits, ".+")
		}
	} else if regSegmentRegexp.MatchString(segment) {
		seg := segment[4 : len(segment)-1]
		splits = strings.SplitN(seg, ">", 2)
	}
	if splits != nil && len(splits) == 2 {
		paramName = splits[0]
		regstr := splits[1]
		if !strings.HasPrefix(regstr, "^") {
			regstr = "^" + regstr
		}
		if !strings.HasSuffix(regstr, "$") {
			regstr = regstr + "$"
		}
		paramReg = regexp.MustCompile(regstr)
	}
	return
}

// AddRoute adds a new route with method, pattern and handler
func (n *Node) AddRoute(method Method, pattern string, h interface{}, filters ...Filter) *Leaf {
	handler := parseHandler(h)
	segments := splitSegments(pattern)
	if _, l, ok := n.add(segments, method, handler, filters...); ok {
		return l
	}
	panic(fmt.Errorf("pattern: %s has been registered", pattern))
}

func (n *Node) add(segments []string, method Method, handler Handler, filters ...Filter) (*Node, *Leaf, bool) {
	if len(segments) == 0 {
		l, ok := n.handle(method, handler, filters...)
		return n, l, ok
	}

	segment := segments[0]
	next := n.findOrCreateNext(segment)
	return next.add(segments[1:], method, handler, filters...)
}

func (n *Node) handle(method Method, handler Handler, filters ...Filter) (l *Leaf, ok bool) {
	if _, ok := n.leaves[method]; ok {
		return nil, false
	}
	l = NewLeaf(n, method, handler)
	n.leaves[method] = l
	l.AddFilters(filters...)
	return l, true
}

func (n *Node) findOrCreateNext(segment string) (next *Node) {
	if !isReg(segment) {
		if ne, ok := n.rawChildren[segment]; ok {
			next = ne
		} else {
			next = NewNode(n.h, segment, n.depth+1)
			n.rawChildren[segment] = next
		}
	} else {
		found := false
		for _, next = range n.regChildren {
			if next.segment == segment {
				found = true
				break
			}
		}
		if !found {
			next = NewNode(n.h, segment, n.depth+1)
			n.regChildren = append(n.regChildren, next)
		}
	}
	next.parent = n
	return
}

func (n *Node) findNext(segment string) (next *Node) {
	if ne, ok := n.rawChildren[segment]; ok {
		next = ne
	} else {
		for _, ne := range n.regChildren {
			if ne.MatchRegexp(segment) {
				next = ne
				break
			}
		}
	}
	return
}

func (n *Node) Serve(ctx *Context) {
	segment := ctx.segment()

	// path matches
	if len(segment) == 0 {
		n.doServe(ctx)
		return
	}

	// find next node
	next := n.findNext(segment)
	if next != nil {
		if next.paramReg != nil {
			ctx.Params()[next.paramName] = segment
		}
		next.Serve(ctx)
		return
	}
	// 404 not found
	ctx.OnError(http.StatusNotFound)
}

func (n *Node) doServe(ctx *Context) {
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
	// ANY matches
	if l, ok := n.leaves["ANY"]; ok {
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

// MatchRegexp checks if the segment match regexp in node
func (n *Node) MatchRegexp(segment string) bool {
	if n.paramReg != nil && n.paramReg.MatchString(segment) {
		return true
	}
	return false
}

// Path returns the full path from root to the node
func (n *Node) Path() string {
	if n.parent == nil {
		return "/"
	}
	ppath := n.parent.Path()
	if ppath == "/" {
		ppath = ""
	}
	if n.paramName != "" {
		return ppath + "/{" + n.paramName + "}"
	}
	return ppath + "/" + n.segment
}

// Parent returns the node's parent node
func (n *Node) Parent() *Node {
	return n.parent
}

// Depth returns the nodes' depth
func (n *Node) Depth() int {
	return n.depth
}

// Segment returns node's segment
func (n *Node) Segment() string {
	return n.segment
}

// Leaves returns all of node's leaves
func (n *Node) Leaves() []*Leaf {
	leaves := make([]*Leaf, len(n.leaves))
	i := 0
	for _, l := range n.leaves {
		leaves[i] = l
		i++
	}
	return leaves[:i]
}

func (n *Node) travel(llist *list.List) {
	for _, l := range n.leaves {
		llist.PushBack(l)
	}
	for _, child := range n.rawChildren {
		child.travel(llist)
	}
	for _, child := range n.regChildren {
		child.travel(llist)
	}
}
