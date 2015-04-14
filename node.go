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

type leaf struct {
	*FilterChain
	h       *Hador
	handler Handler
	method  string
}

func newLeaf(h *Hador, method string, handler Handler) *leaf {
	l := &leaf{
		h:       h,
		method:  method,
		handler: handler,
	}
	l.FilterChain = NewFilterChain(l.handler)
	return l
}

type dispatcher struct {
	h    *Hador
	node *node
}

func (d *dispatcher) Serve(ctx *Context) {
	n := d.node
	segments := genSegments(ctx.Request.URL.Path)
	if len(segments) == n.depth {
		if l, ok := n.leaves[ctx.Request.Method]; ok {
			l.Serve(ctx)
			return
		}
		if l, ok := n.leaves["ANY"]; ok {
			l.Serve(ctx)
			return
		}
		methods := make([]string, len(n.leaves))
		i := 0
		for m := range n.leaves {
			methods[i] = m
			i++
		}
		ctx.MethodNotAllowed(strings.Join(methods, ","))
		return
	}
	segment := segments[n.depth]
	var next *node
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
	ctx.NotFound()
}

type node struct {
	*FilterChain
	h             *Hador
	depth         int
	segment       string
	regexpSegment *regexp.Regexp
	rawChildren   map[string]*node
	regChildren   []*node
	leaves        map[string]*leaf
}

func newNode(h *Hador, segment string, depth int) *node {
	var reg *regexp.Regexp
	if segment != "" && regSegmentRegexp.MatchString(segment) {
		reg = regexp.MustCompile(segment)
	}
	n := &node{
		h:             h,
		segment:       segment,
		depth:         depth,
		regexpSegment: reg,
		rawChildren:   make(map[string]*node),
		regChildren:   make([]*node, 0),
		leaves:        make(map[string]*leaf),
	}
	n.FilterChain = NewFilterChain(&dispatcher{h: n.h, node: n})
	return n
}

func (n *node) Options(pattern string, handler Handler) Beforer {
	return n.AddRoute("OPTIONS", pattern, handler)
}

func (n *node) Get(pattern string, handler Handler) Beforer {
	return n.AddRoute("GET", pattern, handler)
}

func (n *node) Head(pattern string, handler Handler) Beforer {
	return n.AddRoute("HEAD", pattern, handler)
}

func (n *node) Post(pattern string, handler Handler) Beforer {
	return n.AddRoute("POST", pattern, handler)
}

func (n *node) Put(pattern string, handler Handler) Beforer {
	return n.AddRoute("PUT", pattern, handler)
}

func (n *node) Delete(pattern string, handler Handler) Beforer {
	return n.AddRoute("DELETE", pattern, handler)
}

func (n *node) Trace(pattern string, handler Handler) Beforer {
	return n.AddRoute("TRACE", pattern, handler)
}

func (n *node) Connect(pattern string, handler Handler) Beforer {
	return n.AddRoute("CONNECT", pattern, handler)
}

func (n *node) Patch(pattern string, handler Handler) Beforer {
	return n.AddRoute("PATCH", pattern, handler)
}

func (n *node) Any(pattern string, handler Handler) Beforer {
	return n.AddRoute("ANY", pattern, handler)
}

func (n *node) Group(pattern string, f func(Router)) Beforer {
	segments := genSegments(pattern)
	r, _, _ := n.add(segments, "", nil)
	f(r)
	return r
}

func (n *node) AddRoute(method, pattern string, handler Handler) Beforer {
	segments := genSegments(pattern)
	if _, l, ok := n.add(segments, method, handler); ok {
		return l
	}
	panic(fmt.Errorf("pattern: %s has been registered", pattern))
}

func (n *node) add(segments []string, method string, handler Handler) (*node, *leaf, bool) {
	if len(segments) == 0 {
		if method != "" && handler != nil {
			if _, ok := n.leaves[method]; ok {
				return n, nil, false
			}
			l := newLeaf(n.h, method, handler)
			n.leaves[method] = l
			return n, l, true
		}
		return n, nil, true
	}
	segment := segments[0]
	var next *node
	if !regSegmentRegexp.MatchString(segment) {
		if ne, ok := n.rawChildren[segment]; ok {
			next = ne
		} else {
			next = newNode(n.h, segment, n.depth+1)
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
			next = newNode(n.h, segment, n.depth+1)
			n.regChildren = append(n.regChildren, next)
		}
	}
	return next.add(segments[1:], method, handler)
}

func (n *node) MatchRegexp(segment string) (string, string, bool) {
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
