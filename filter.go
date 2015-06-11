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

// CombineFilters combines multi Filters into a single Filter.
func CombineFilters(filters ...Filter) Filter {
	// if no Filter or Filters is nil return nil
	if filters == nil || len(filters) == 0 {
		return nil
	}
	// if the first Filter is the only, return itself
	first := filters[0]
	if len(filters) == 1 {
		return first
	}
	// combine others as second Filter
	second := CombineFilters(filters[1:]...)

	return Combine2Filters(first, second)
}

// Combine2Filters combines 2 Filters into a single Filter.
func Combine2Filters(first, second Filter) Filter {
	// return the other if one is nil
	// return nil if both are nil
	if first == nil {
		return second
	}
	if second == nil {
		return first
	}
	return FilterFunc(func(ctx *Context, next Handler) {
		first.Filter(ctx, filterHandler(second, next))
	})
}

func filterHandler(filter Filter, handler Handler) Handler {
	if filter == nil {
		return handler
	}
	return HandlerFunc(func(ctx *Context) {
		filter.Filter(ctx, handler)
	})
}

// FilterChain struct combines multi Filters and Handler into a single Handler.
type FilterChain struct {
	handler Handler
	filter  Filter
}

// NewFilterChain creates new FilterChain instance
func NewFilterChain(handler Handler, filters ...Filter) *FilterChain {
	if handler == nil {
		panic("handler shouldn't be nil")
	}
	return &FilterChain{
		handler: handler,
		filter:  CombineFilters(filters...),
	}
}

// Serve implements Handler interface
func (fc *FilterChain) Serve(ctx *Context) {
	filterHandler(fc.filter, fc.handler).Serve(ctx)
}

// Before implements Beforer interface
func (fc *FilterChain) Before(filter Filter) Beforer {
	fc.filter = Combine2Filters(fc.filter, filter)
	return fc
}

// BeforeFunc implements Beforer interface
func (fc *FilterChain) BeforeFunc(f func(*Context, Handler)) Beforer {
	return fc.Before(FilterFunc(f))
}

// AddFilters adds all filters to chain
func (fc *FilterChain) AddFilters(filters ...Filter) {
	fc.filter = Combine2Filters(fc.filter, CombineFilters(filters...))
}

// InsertAfter insert filters before self
func (fc *FilterChain) InsertFront(filters ...Filter) {
	fc.filter = Combine2Filters(CombineFilters(filters...), fc.filter)
}

// InsertBack inserts filters after self
func (fc *FilterChain) InsertBack(filters ...Filter) {
	fc.filter = Combine2Filters(fc.filter, CombineFilters(filters...))
}
