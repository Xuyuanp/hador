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

var regSegmentRegexp = regexp.MustCompile(`\(\?P<.+>.+\)`)

// Leaf struct
type Leaf struct {
	*FilterChain
	h       *Hador
	Parent  *Node
	Handler Handler
	Method  string
}

// NewLeaf creates new Leaf instance
func NewLeaf(h *Hador, method string, handler Handler) *Leaf {
	l := &Leaf{
		h:       h,
		Method:  method,
		Handler: handler,
	}
	l.FilterChain = NewFilterChain(l.Handler)
	return l
}

type dispatcher struct {
	h    *Hador
	node *Node
}

func (d *dispatcher) Serve(ctx *Context) {
	n := d.node
	segments := ctx.Request.segments
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
			if key, value, ok := ne.MatchRegexp(segment); ok && key != "" {
				ctx.Params[key] = value
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
	h             *Hador
	Parent        *Node
	Depth         int
	Segment       string
	regexpSegment *regexp.Regexp
	rawChildren   map[string]*Node
	regChildren   []*Node
	Leaves        map[string]*Leaf
}

// NewNode creates new Node instance.
func NewNode(h *Hador, segment string, depth int) *Node {
	var reg *regexp.Regexp
	if segment != "" && regSegmentRegexp.MatchString(segment) {
		reg = regexp.MustCompile(segment)
	}
	n := &Node{
		h:             h,
		Segment:       segment,
		Depth:         depth,
		regexpSegment: reg,
		rawChildren:   make(map[string]*Node),
		regChildren:   make([]*Node, 0),
		Leaves:        make(map[string]*Leaf),
	}
	n.FilterChain = NewFilterChain(&dispatcher{h: n.h, node: n})
	return n
}

// Options adds route by call AddRoute
func (n *Node) Options(pattern string, handler Handler) Beforer {
	return n.AddRoute("OPTIONS", pattern, handler)
}

// Get adds route by call AddRoute
func (n *Node) Get(pattern string, handler Handler) Beforer {
	return n.AddRoute("GET", pattern, handler)
}

// Head adds route by call AddRoute
func (n *Node) Head(pattern string, handler Handler) Beforer {
	return n.AddRoute("HEAD", pattern, handler)
}

// Post adds route by call AddRoute
func (n *Node) Post(pattern string, handler Handler) Beforer {
	return n.AddRoute("POST", pattern, handler)
}

// Put adds route by call AddRoute
func (n *Node) Put(pattern string, handler Handler) Beforer {
	return n.AddRoute("PUT", pattern, handler)
}

// Delete adds route by call AddRoute
func (n *Node) Delete(pattern string, handler Handler) Beforer {
	return n.AddRoute("DELETE", pattern, handler)
}

// Trace adds route by call AddRoute
func (n *Node) Trace(pattern string, handler Handler) Beforer {
	return n.AddRoute("TRACE", pattern, handler)
}

// Connect adds route by call AddRoute
func (n *Node) Connect(pattern string, handler Handler) Beforer {
	return n.AddRoute("CONNECT", pattern, handler)
}

// Patch adds route by call AddRoute
func (n *Node) Patch(pattern string, handler Handler) Beforer {
	return n.AddRoute("PATCH", pattern, handler)
}

// Any adds route by call AddRoute
func (n *Node) Any(pattern string, handler Handler) Beforer {
	return n.AddRoute("ANY", pattern, handler)
}

// Group adds group routes
func (n *Node) Group(pattern string, f func(Router)) Beforer {
	segments := genSegments(pattern)
	r, _, _ := n.add(segments, "", nil)
	f(r)
	return r
}

// AddRoute adds a new route with method, pattern and handler
func (n *Node) AddRoute(method, pattern string, handler Handler) Beforer {
	segments := genSegments(pattern)
	if _, l, ok := n.add(segments, method, handler); ok {
		return l
	}
	panic(fmt.Errorf("pattern: %s has been registered", pattern))
}

func (n *Node) add(segments []string, method string, handler Handler) (*Node, *Leaf, bool) {
	if len(segments) == 0 {
		if method != "" && handler != nil {
			if _, ok := n.Leaves[method]; ok {
				return n, nil, false
			}
			l := NewLeaf(n.h, method, handler)
			l.Parent = n
			n.Leaves[method] = l
			return n, l, true
		}
		return n, nil, true
	}
	segment := segments[0]
	var next *Node
	if !regSegmentRegexp.MatchString(segment) {
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
	return next.add(segments[1:], method, handler)
}

// MatchRegexp checks if the segment match regexp in node
func (n *Node) MatchRegexp(segment string) (string, string, bool) {
	if n.regexpSegment != nil {
		result := n.regexpSegment.FindStringSubmatch(segment)
		if len(result) > 1 && result[0] == segment {
			value := result[1]
			key := n.regexpSegment.SubexpNames()[1]
			return key, value, true
		}
	}
	return "", "", false
}
