package hador

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func runRequest(b *testing.B, h *Hador, method, path string) {
	req, _ := http.NewRequest(method, path, nil)
	resp := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp.Body.Reset()
		h.ServeHTTP(resp, req)
	}
}

func BenchmarkLeaf(b *testing.B) {
	n := &node{}
	l := n.AddRoute(GET, "/ping", func(_ *Context) {})
	ctx := newContext(defaultLogger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		l.Serve(ctx)
	}
}

func BenchmarkMatchLeaf(b *testing.B) {
	n := &node{}
	l := n.AddRoute(GET, "/ping", func(_ *Context) {})
	parent := l.parent

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parent.matchLeaf(GET)
	}
}

func BenchmarkNodeMatch(b *testing.B) {
	n := &node{}
	n.AddRoute(GET, "/hello/{name}", func(*Context) {})

	ctx := newContext(defaultLogger)
	ctx.params = make(Params, n.findMaxParams())
	ctx.params = ctx.params[0:0]
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx.params = ctx.params[0:0]
		_, l, _ := n.match(GET, "/hello/jack", ctx.Params())
		l.Serve(ctx)
	}
}

func BenchmarkHadorServe(b *testing.B) {
	h := New()
	h.Get("/ping", func(*Context) {})

	req, _ := http.NewRequest("GET", "/ping", nil)
	recorder := httptest.NewRecorder()
	resp := NewResponseWriter(recorder)

	ctx := newContext(h.Logger)
	ctx.params = make(Params, h.root.findMaxParams())

	ctx.reset(resp, req)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx.params = ctx.params[0:0]
		h.Serve(ctx)
	}
}

func BenchmarkSinglePath(b *testing.B) {
	h := New()
	h.Get("/ping", func(ctx *Context) {})

	runRequest(b, h, "GET", "/ping")
}

func BenchmarkParam(b *testing.B) {
	h := New()
	h.Get("/{name}", func(ctx *Context) {})

	runRequest(b, h, "GET", "/jack")
}

func Benchmark5Param(b *testing.B) {
	h := New()
	h.Get("/{a}/{b}/{c}/{d}/{e}", func(ctx *Context) {})

	runRequest(b, h, "GET", "/a/b/c/d/e")
}

func BenchmarkCombine2Filters(b *testing.B) {
	f1 := FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) })
	f2 := FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) })
	handler := HandlerFunc(func(*Context) {})

	f3 := Combine2Filters(f1, f2)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		f3.Filter(nil, handler)
	}
}

func BenchmarkFilterChain(b *testing.B) {
	handler := HandlerFunc(func(*Context) {})

	f1 := FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) })
	f2 := FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) })
	fc := NewFilterChain(handler, f1)
	fc.AddFilters(f2)

	b.ResetTimer()
	b.ReportAllocs()

	ctx := newContext(defaultLogger)

	for i := 0; i < b.N; i++ {
		fc.Serve(ctx)
	}
}

func BenchmarkGroup(b *testing.B) {
	h := New()
	h.Group("/group", func(router Router) {
		router.Get("/ping", HandlerFunc(func(ctx *Context) {}))
	})

	runRequest(b, h, "GET", "/group/ping")
}

func BenchmarkGroupFilters(b *testing.B) {
	h := New()
	h.AddFilters(
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
	)
	h.Group("/group", func(router Router) {
		router.Get("/ping", HandlerFunc(func(ctx *Context) {}),
			FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
			FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		)
	},
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
	)

	runRequest(b, h, "GET", "/group/ping")
}

func BenchmarkHTTPHandler(b *testing.B) {
	h := New()
	h.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})
	runRequest(b, h, "GET", "/ping")
}

func BenchmarkController(b *testing.B) {
	h := New()
	h.AddController("/controller", &testController{prepared: true})
	runRequest(b, h, "GET", "/controller")
}

func BenchmarkRenderJSON(b *testing.B) {
	h := New()
	h.Get("/ping", func(ctx *Context) {
		ctx.RenderJSON(struct {
			Status int
			Body   string
		}{
			Status: 200,
			Body:   "OK",
		})
	})

	runRequest(b, h, "GET", "/ping")
}

func BenchmarkEmptyServeHTTP(b *testing.B) {
	h := New()

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		h._ServeHTTP(resp, req)
	}
}
