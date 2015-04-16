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

type dispatcher struct {
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
	Parent        *Node
	Depth         int
	Segment       string
	regexpSegment *regexp.Regexp
	rawChildren   map[string]*Node
	regChildren   []*Node
	Leaves        map[string]*Leaf
}

// NewNode creates new Node instance.
func NewNode(segment string, depth int) *Node {
	var reg *regexp.Regexp
	if segment != "" && regSegmentRegexp.MatchString(segment) {
		reg = regexp.MustCompile(segment)
	}
	n := &Node{
		Segment:       segment,
		Depth:         depth,
		regexpSegment: reg,
		rawChildren:   make(map[string]*Node),
		regChildren:   make([]*Node, 0),
		Leaves:        make(map[string]*Leaf),
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
			l := NewLeaf(method, handler)
			l.Parent = n
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
	if !regSegmentRegexp.MatchString(segment) {
		if ne, ok := n.rawChildren[segment]; ok {
			next = ne
		} else {
			next = NewNode(segment, n.Depth+1)
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
			next = NewNode(segment, n.Depth+1)
			n.regChildren = append(n.regChildren, next)
		}
	}
	next.Parent = n
	return next.add(segments[1:], method, handler, filters...)
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

// Path returns the full path from root to the node
func (n *Node) Path() string {
	if n.Parent == nil {
		return "/"
	}
	ppath := n.Parent.Path()
	if ppath == "/" {
		ppath = ""
	}
	return ppath + "/" + n.Segment
}
