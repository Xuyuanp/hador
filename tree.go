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

type dispatcher struct {
	tree *tree
}

func (d dispatcher) Serve(ctx *Context) {
	t := d.tree
	segments := genSegments(ctx.Request.URL.Path)
	if len(segments) == t.depth {
		t.prvFilterChain.Serve(ctx)
		return
	}
	segment := segments[t.depth]
	var next *tree
	if n, ok := t.rawChildren[segment]; ok {
		next = n
	} else {
		for _, n := range t.regChildren {
			if key, value, ok := n.MatchRegexp(segment); ok && key != "" {
				ctx.Params[key] = value
				next = n
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

type tree struct {
	depth          int
	segment        string
	regexpSegment  *regexp.Regexp
	pubFilterChain *FilterChain
	prvFilterChain *FilterChain
	handler        *MethodHandler
	rawChildren    map[string]*tree
	regChildren    []*tree
}

var regSegmentRegexp = regexp.MustCompile(`\(\?P<.+>.+\)`)

func newTree(segment string, depth int) *tree {
	var reg *regexp.Regexp
	if segment != "" && regSegmentRegexp.MatchString(segment) {
		reg = regexp.MustCompile(segment)
	}
	t := &tree{
		depth:         depth,
		segment:       segment,
		regexpSegment: reg,
		handler:       NewMethodHandler(),
		rawChildren:   make(map[string]*tree),
		regChildren:   make([]*tree, 0),
	}
	t.prvFilterChain = NewFilterChain(t.handler)
	t.pubFilterChain = NewFilterChain(&dispatcher{tree: t})
	return t
}

func (t *tree) Options(pattern string, handler Handler) Beforer {
	return t.AddRoute("OPTIONS", pattern, handler)
}

func (t *tree) Get(pattern string, handler Handler) Beforer {
	return t.AddRoute("GET", pattern, handler)
}

func (t *tree) Head(pattern string, handler Handler) Beforer {
	return t.AddRoute("HEAD", pattern, handler)
}

func (t *tree) Post(pattern string, handler Handler) Beforer {
	return t.AddRoute("POST", pattern, handler)
}

func (t *tree) Put(pattern string, handler Handler) Beforer {
	return t.AddRoute("PUT", pattern, handler)
}

func (t *tree) Delete(pattern string, handler Handler) Beforer {
	return t.AddRoute("DELETE", pattern, handler)
}

func (t *tree) Trace(pattern string, handler Handler) Beforer {
	return t.AddRoute("TRACE", pattern, handler)
}

func (t *tree) Connect(pattern string, handler Handler) Beforer {
	return t.AddRoute("CONNECT", pattern, handler)
}

func (t *tree) Patch(pattern string, handler Handler) Beforer {
	return t.AddRoute("PATCH", pattern, handler)
}

func (t *tree) Any(pattern string, handler Handler) Beforer {
	return t.AddRoute("ANY", pattern, handler)
}

func (t *tree) Group(pattern string, f func(Router)) Beforer {
	segments := genSegments(pattern)
	r, _ := t.add(segments, "", nil)
	f(r)
	return t.pubFilterChain
}

func (t *tree) AddRoute(method string, pattern string, handler Handler) Beforer {
	segments := genSegments(pattern)
	if b, ok := t.add(segments, method, handler); ok {
		return b.prvFilterChain
	}
	panic(fmt.Errorf("pattern: %s has been registered", pattern))
}

func (t *tree) add(segments []string, method string, handler Handler) (*tree, bool) {
	if len(segments) == 0 {
		if method != "" && handler != nil {
			ok := t.handler.Handle(method, handler)
			return t, ok
		}
		return t, true
	}
	segment := segments[0]
	var next *tree
	if !regSegmentRegexp.MatchString(segment) {
		if n, ok := t.rawChildren[segment]; ok {
			next = n
		} else {
			next = newTree(segment, t.depth+1)
			t.rawChildren[segment] = next
		}
	} else {
		found := false
		for _, next := range t.regChildren {
			if next.segment == segment {
				found = true
				break
			}
		}
		if !found {
			next = newTree(segment, t.depth+1)
			t.regChildren = append(t.regChildren, next)
		}
	}
	return next.add(segments[1:], method, handler)
}

func (t *tree) MatchRegexp(segment string) (string, string, bool) {
	if t.regexpSegment != nil {
		result := t.regexpSegment.FindStringSubmatch(segment)
		if len(result) > 1 && result[0] == segment {
			value := result[1]
			key := t.regexpSegment.SubexpNames()[1]
			return key, value, true
		}
	}
	return "", "", false
}

func (t *tree) Serve(ctx *Context) {
	t.pubFilterChain.Serve(ctx)
}
