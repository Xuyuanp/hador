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

// Filter interface
type Filter interface {
	Filter(ctx *Context, next Handler)
}

// FilterFunc as function
type FilterFunc func(*Context, Handler)

// Filter implements Filter interface
func (ff FilterFunc) Filter(ctx *Context, next Handler) {
	ff(ctx, next)
}

// Beforer interface
type Beforer interface {
	Before(Filter) Beforer
	BeforeFunc(func(*Context, Handler)) Beforer
}

// FilterChain struct
type FilterChain struct {
	handler Handler
	filter  Filter
	next    *FilterChain
}

// NewFilterChain creates new FilterChain instance
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

// Serve implements Handler interface
func (fc *FilterChain) Serve(ctx *Context) {
	if fc.handler != nil {
		fc.handler.Serve(ctx)
	} else {
		fc.filter.Filter(ctx, fc.next)
	}
}

// Before implements Beforer interface
func (fc *FilterChain) Before(filter Filter) Beforer {
	if filter == nil {
		return fc
	}
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

// BeforeFunc implements Beforer interface
func (fc *FilterChain) BeforeFunc(f func(*Context, Handler)) Beforer {
	return fc.Before(FilterFunc(f))
}

// AddFilters adds all filters to chain
func (fc *FilterChain) AddFilters(filters ...Filter) {
	for _, filter := range filters {
		fc.Before(filter)
	}
}
