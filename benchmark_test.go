package hador

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func runRequest(b *testing.B, h *Hador, method, path string) {
	req, _ := http.NewRequest(method, "", nil)
	req.URL.Path = path
	resp := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(resp, req)
	}
}

func BenchmarkSinglePath(b *testing.B) {
	h := New()
	h.Get("/ping", HandlerFunc(func(ctx *Context) {}))

	runRequest(b, h, "GET", "/ping")
}

func BenchmarkParam(b *testing.B) {
	h := New()
	h.Get("/{name}", HandlerFunc(func(ctx *Context) {}))

	runRequest(b, h, "GET", "/jack")
}

func Benchmark5Param(b *testing.B) {
	h := New()
	h.Get("/{a}/{b}/{c}/{d}/{e}", HandlerFunc(func(ctx *Context) {}))

	runRequest(b, h, "GET", "/a/b/c/d/e")
}

func BenchmarkFilters(b *testing.B) {
	h := New()
	h.Get("/ping", HandlerFunc(func(ctx *Context) {}),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
		FilterFunc(func(ctx *Context, next Handler) { next.Serve(ctx) }),
	)

	runRequest(b, h, "GET", "/ping")
}

func BenchmarkRenderJSON(b *testing.B) {
	h := New()
	h.Get("/ping", HandlerFunc(func(ctx *Context) {
		ctx.RenderJSON(struct {
			Status int    `json:"status"`
			Body   string `json:"body"`
		}{
			Status: 200,
			Body:   "OK",
		})
	}))

	runRequest(b, h, "GET", "/ping")
}
