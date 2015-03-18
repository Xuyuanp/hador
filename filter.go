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

type Filter interface {
	Filter(ctx *Context, next Handler)
}

type FilterFunc func(*Context, Handler)

func (ff FilterFunc) Filter(ctx *Context, next Handler) {
	ff(ctx, next)
}

type Beforer interface {
	Before(Filter) Beforer
	BeforeFunc(func(*Context, Handler)) Beforer
}

type FilterChain struct {
	handler Handler
	filter  Filter
	next    *FilterChain
}

func NewFilterChain(handler Handler, filters ...Filter) *FilterChain {
	if handler == nil {
		panic("handler shouldn't be nil")
	}
	if len(filters) == 0 {
		return &FilterChain{
			handler: handler,
		}
	}
	filter := filters[0]
	if filter == nil {
		panic("filter shouldn't be nil")
	}
	return &FilterChain{
		filter: filter,
		next:   NewFilterChain(handler, filters[1:]...),
	}
}

func (fc *FilterChain) ServeHTTP(ctx *Context) {
	if fc.handler != nil {
		fc.handler.ServeHTTP(ctx)
	} else {
		fc.filter.Filter(ctx, fc.next)
	}
}

func (fc *FilterChain) Before(filter Filter) Beforer {
	tmp := fc
	for tmp.next != nil {
		tmp = tmp.next
	}
	tmp.next = &FilterChain{
		handler: tmp.handler,
	}
	tmp.handler = nil
	tmp.filter = filter
	return fc
}

func (fc *FilterChain) BeforeFunc(f func(*Context, Handler)) Beforer {
	return fc.Before(FilterFunc(f))
}
