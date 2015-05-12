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
	"regexp"
	"strings"
)

func genSegments(path string) []string {
	if path == "/" {
		path = ""
	}
	return strings.Split(path, "/")[1:]
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

type dispatcher struct {
	node *Node
}

func (d *dispatcher) Serve(ctx *Context) {
	n := d.node
	segments := ctx.segments
	// path matches
	if len(segments) == n.Depth {
		// 404 not found
		if len(n.Leaves) == 0 {
			ctx.NotFound()
			return
		}
		// method matches
		if l, ok := n.Leaves[ctx.Request.Method]; ok {
			l.Serve(ctx)
			return
		}
		// ANY matches
		if l, ok := n.Leaves["ANY"]; ok {
			l.Serve(ctx)
			return
		}
		// 405 method not allowed
		methods := make([]string, len(n.Leaves))
		i := 0
		for m := range n.Leaves {
			methods[i] = m
			i++
		}
		ctx.MethodNotAllowed(methods)
		return
	}
	// find next node
	segment := segments[n.Depth]
	var next *Node
	if ne, ok := n.rawChildren[segment]; ok {
		next = ne
	} else {
		for _, ne := range n.regChildren {
			if ne.MatchRegexp(segment) {
				ctx.Params[ne.paramName] = segment
				next = ne
				break
			}
		}
	}
	if next != nil {
		next.Serve(ctx)
		return
	}
	// 404 not found
	ctx.NotFound()
}

// Node struct
type Node struct {
	*FilterChain
	h           *Hador
	Parent      *Node
	Depth       int
	Segment     string
	paramName   string
	paramReg    *regexp.Regexp
	rawChildren map[string]*Node
	regChildren []*Node
	Leaves      map[string]*Leaf
}

// NewNode creates new Node instance.
func NewNode(h *Hador, segment string, depth int) *Node {
	var paramName string
	var paramReg *regexp.Regexp
	if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
		seg := segment[1 : len(segment)-1]
		splits := strings.SplitN(seg, ":", 2)
		if len(splits) == 1 {
			splits = append(splits, ".+")
		}
		paramName = splits[0]
		paramReg = regexp.MustCompile(splits[1])
	} else if regSegmentRegexp.MatchString(segment) {
		seg := segment[4 : len(segment)-1]
		splits := strings.SplitN(seg, ">", 2)
		paramName = splits[0]
		paramReg = regexp.MustCompile(splits[1])
	}
	n := &Node{
		h:           h,
		Segment:     segment,
		Depth:       depth,
		paramName:   paramName,
		paramReg:    paramReg,
		rawChildren: make(map[string]*Node),
		regChildren: make([]*Node, 0),
		Leaves:      make(map[string]*Leaf),
	}
	n.FilterChain = NewFilterChain(&dispatcher{node: n})
	return n
}

// Options adds route by call AddRoute
func (n *Node) Options(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("OPTIONS", pattern, handler, filters...)
}

// Get adds route by call AddRoute
func (n *Node) Get(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("GET", pattern, handler, filters...)
}

// Head adds route by call AddRoute
func (n *Node) Head(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("HEAD", pattern, handler, filters...)
}

// Post adds route by call AddRoute
func (n *Node) Post(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("POST", pattern, handler, filters...)
}

// Put adds route by call AddRoute
func (n *Node) Put(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("PUT", pattern, handler, filters...)
}

// Delete adds route by call AddRoute
func (n *Node) Delete(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("DELETE", pattern, handler, filters...)
}

// Trace adds route by call AddRoute
func (n *Node) Trace(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("TRACE", pattern, handler, filters...)
}

// Connect adds route by call AddRoute
func (n *Node) Connect(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("CONNECT", pattern, handler, filters...)
}

// Patch adds route by call AddRoute
func (n *Node) Patch(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("PATCH", pattern, handler, filters...)
}

// Any adds route by call AddRoute
func (n *Node) Any(pattern string, handler Handler, filters ...Filter) *Leaf {
	return n.AddRoute("ANY", pattern, handler, filters...)
}

// Group adds group routes
func (n *Node) Group(pattern string, f func(Router), filters ...Filter) {
	segments := genSegments(pattern)
	r, _, _ := n.add(segments, "", nil, filters...)
	f(r)
}

// AddRoute adds a new route with method, pattern and handler
func (n *Node) AddRoute(method, pattern string, handler Handler, filters ...Filter) *Leaf {
	segments := genSegments(pattern)
	if _, l, ok := n.add(segments, method, handler, filters...); ok {
		return l
	}
	panic(fmt.Errorf("pattern: %s has been registered", pattern))
}

func (n *Node) add(segments []string, method string, handler Handler, filters ...Filter) (*Node, *Leaf, bool) {
	if len(segments) == 0 {
		if method != "" && handler != nil {
			if _, ok := n.Leaves[method]; ok {
				return n, nil, false
			}
			l := NewLeaf(n.h, method, handler)
			l.parent = n
			l.path = n.Path()
			n.Leaves[method] = l
			l.AddFilters(filters...)
			return n, l, true
		}
		n.AddFilters(filters...)
		return n, nil, true
	}
	segment := segments[0]
	var next *Node
	if !isReg(segment) {
		if ne, ok := n.rawChildren[segment]; ok {
			next = ne
		} else {
			next = NewNode(n.h, segment, n.Depth+1)
			n.rawChildren[segment] = next
		}
	} else {
		found := false
		for _, next = range n.regChildren {
			if next.Segment == segment {
				found = true
				break
			}
		}
		if !found {
			next = NewNode(n.h, segment, n.Depth+1)
			n.regChildren = append(n.regChildren, next)
		}
	}
	next.Parent = n
	return next.add(segments[1:], method, handler, filters...)
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
	if n.Parent == nil {
		return "/"
	}
	ppath := n.Parent.Path()
	if ppath == "/" {
		ppath = ""
	}
	if n.paramName != "" {
		return ppath + "/{" + n.paramName + "}"
	}
	return ppath + "/" + n.Segment
}

func (n *Node) travel(llist *list.List) {
	for _, l := range n.Leaves {
		llist.PushBack(l)
	}
	for _, child := range n.rawChildren {
		child.travel(llist)
	}
	for _, child := range n.regChildren {
		child.travel(llist)
	}
}
