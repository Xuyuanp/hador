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

func BenchmarkFilters(b *testing.B) {
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
	h.Any("/controller", &testController{prepared: true})
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
