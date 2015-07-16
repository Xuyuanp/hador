package hador

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGinNode(t *testing.T) {
	convey.Convey("Test gin_node", t, func() {
		convey.Convey("Test AddRoute", func() {
			convey.Convey("Test static", func() {
				convey.Convey("single route", func() {
					n := &node{}
					n.AddRoute(GET, "/init", func(*Context) {})
					convey.So(n.segment, convey.ShouldEqual, "/init")
					convey.So(n.ntype, convey.ShouldEqual, static)
				})
				convey.Convey("two routes", func() {
					n := &node{}
					n.AddRoute(GET, "/initfoo", func(*Context) {})
					n.AddRoute(GET, "/initbar", func(*Context) {})
					convey.So(n.segment, convey.ShouldEqual, "/init")
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.indices, convey.ShouldEqual, "fb")
					convey.So(len(n.children), convey.ShouldEqual, 2)
					convey.So(n.children[0].segment, convey.ShouldEqual, "foo")
					convey.So(n.children[1].segment, convey.ShouldEqual, "bar")
				})
				convey.Convey("two routes, longer first", func() {
					n := &node{}
					n.AddRoute(GET, "/initfoo", func(*Context) {})
					n.AddRoute(GET, "/init", func(*Context) {})
					convey.So(n.segment, convey.ShouldEqual, "/init")
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.indices, convey.ShouldEqual, "f")
					convey.So(len(n.children), convey.ShouldEqual, 1)
					convey.So(n.children[0].segment, convey.ShouldEqual, "foo")
				})
				convey.Convey("two routes, shorter first", func() {
					n := &node{}
					n.AddRoute(GET, "/init", func(*Context) {})
					n.AddRoute(GET, "/initfoo", func(*Context) {})
					convey.So(n.segment, convey.ShouldEqual, "/init")
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.indices, convey.ShouldEqual, "f")
					convey.So(len(n.children), convey.ShouldEqual, 1)
					convey.So(n.children[0].segment, convey.ShouldEqual, "foo")
				})
				convey.Convey("two routes,  only the same slash", func() {
					n := &node{}
					n.AddRoute(GET, "/foo", func(*Context) {})
					n.AddRoute(GET, "/bar", func(*Context) {})
					convey.So(n.segment, convey.ShouldEqual, "/")
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.indices, convey.ShouldEqual, "fb")
					convey.So(len(n.children), convey.ShouldEqual, 2)
					convey.So(n.children[0].segment, convey.ShouldEqual, "foo")
					convey.So(n.children[1].segment, convey.ShouldEqual, "bar")
				})
			})
			convey.Convey("Test param", func() {
				convey.Convey("param at root", func() {
					n := &node{}
					n.AddRoute(GET, "/{name}", func(*Context) {})
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.segment, convey.ShouldEqual, "/")
					convey.So(n.paramChild, convey.ShouldNotBeNil)
					convey.So(n.paramChild.ntype, convey.ShouldEqual, param)
					convey.So(n.paramChild.segment, convey.ShouldEqual, "{name}")
				})
				convey.Convey("param not at root", func() {
					n := &node{}
					n.AddRoute(GET, "/foo/{name}", func(*Context) {})
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.segment, convey.ShouldEqual, "/foo/")
					convey.So(n.paramChild, convey.ShouldNotBeNil)
					convey.So(n.paramChild.ntype, convey.ShouldEqual, param)
					convey.So(n.paramChild.segment, convey.ShouldEqual, "{name}")
				})
				convey.Convey("param between statics", func() {
					n := &node{}
					n.AddRoute(GET, "/foo/{name}/bar", func(*Context) {})
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.segment, convey.ShouldEqual, "/foo/")
					convey.So(n.paramChild, convey.ShouldNotBeNil)
					convey.So(n.paramChild.ntype, convey.ShouldEqual, param)
					convey.So(n.paramChild.segment, convey.ShouldEqual, "{name}")
					convey.So(n.paramChild.indices, convey.ShouldEqual, "/")
					convey.So(len(n.paramChild.children), convey.ShouldEqual, 1)
					convey.So(n.paramChild.children[0].ntype, convey.ShouldEqual, static)
					convey.So(n.paramChild.children[0].segment, convey.ShouldEqual, "/bar")
				})
				convey.Convey("multi param nodes", func() {
					n := &node{}
					n.AddRoute(GET, "/foo/{name}/bar", func(*Context) {})
					n.AddRoute(GET, "/foo/{name}/fizz", func(*Context) {})
					convey.So(n.ntype, convey.ShouldEqual, static)
					convey.So(n.segment, convey.ShouldEqual, "/foo/")
					convey.So(n.paramChild, convey.ShouldNotBeNil)
					convey.So(n.paramChild.ntype, convey.ShouldEqual, param)
					convey.So(n.paramChild.segment, convey.ShouldEqual, "{name}")
					convey.So(n.paramChild.indices, convey.ShouldEqual, "/")
					convey.So(len(n.paramChild.children), convey.ShouldEqual, 1)
					convey.So(n.paramChild.children[0].ntype, convey.ShouldEqual, static)
					convey.So(n.paramChild.children[0].segment, convey.ShouldEqual, "/")
					convey.So(n.paramChild.children[0].indices, convey.ShouldEqual, "bf")
					convey.So(n.paramChild.children[0].children[0].segment, convey.ShouldEqual, "bar")
					convey.So(n.paramChild.children[0].children[1].segment, convey.ShouldEqual, "fizz")
				})
				convey.Convey("Test fuck", func() {
					convey.Convey("test users", func() {
						n := &node{}
						n.AddRoute(GET, "/hello/{name}", func(*Context) {})
						n.AddRoute(POST, "/hello/{name}", func(*Context) {})
						n.AddRoute(GET, "/hello/{name}/today", func(*Context) {})
						n.AddRoute(POST, "/hello/{name}/today", func(*Context) {})
						convey.So(n.ntype, convey.ShouldEqual, static)
						convey.So(n.segment, convey.ShouldEqual, "/hello/")
						convey.So(n.paramChild, convey.ShouldNotBeNil)
						child := n.paramChild
						convey.So(child.ntype, convey.ShouldEqual, param)
						convey.So(child.segment, convey.ShouldEqual, "{name}")
						convey.So(child.ntype, convey.ShouldEqual, param)
						convey.So(len(child.leaves), convey.ShouldEqual, 2)
						convey.So(child.indices, convey.ShouldEqual, "/")
						convey.So(len(child.children), convey.ShouldEqual, 1)
						child = child.children[0]
						convey.So(child.ntype, convey.ShouldEqual, static)
						convey.So(child.segment, convey.ShouldEqual, "/today")
						convey.So(child.leaves, convey.ShouldNotBeNil)
					})
				})
			})
		})
		convey.Convey("Test find", func() {
			convey.Convey("Test static", func() {
				convey.Convey("root", func() {
					n := &node{}
					l := n.AddRoute(GET, "/", func(*Context) {})
					convey.So(n.find(GET, "/"), convey.ShouldEqual, l)
				})
				convey.Convey("single route", func() {
					n := &node{}
					l := n.AddRoute(GET, "/foo", func(*Context) {})
					convey.So(n.find(GET, "/foo"), convey.ShouldEqual, l)
				})
				convey.Convey("multi route with same prefix", func() {
					n := &node{}
					l1 := n.AddRoute(GET, "/initfoo", func(*Context) {})
					l2 := n.AddRoute(GET, "/initbar", func(*Context) {})
					convey.So(n.find(GET, "/initfoo"), convey.ShouldEqual, l1)
					convey.So(n.find(GET, "/initbar"), convey.ShouldEqual, l2)
				})
				convey.Convey("multi routes without same prefix", func() {
					n := &node{}
					n.AddRoute(GET, "/foo", func(*Context) {})
					n.AddRoute(GET, "/bar", func(*Context) {})
					convey.So(n.find(GET, "/foo"), convey.ShouldNotBeNil)
				})
				convey.Convey("multi long routes", func() {
					n := &node{}
					l1 := n.AddRoute(GET, "/abc/def/ghi/jkl/mno/pqr", func(*Context) {})
					l2 := n.AddRoute(GET, "/abc/xyz/ghi/jkl/mno/pqr", func(*Context) {})
					convey.So(n.find(GET, "/abc/def/ghi/jkl/mno/pqr"), convey.ShouldEqual, l1)
					convey.So(n.find(GET, "/abc/xyz/ghi/jkl/mno/pqr"), convey.ShouldEqual, l2)
				})
			})
		})

		convey.Convey("Test param route", func() {
			convey.Convey("/param route", func() {
				n := &node{}
				l := n.AddRoute(GET, "/{name}", func(*Context) {})
				convey.So(n.find(GET, "/foo"), convey.ShouldEqual, l)
			})
			convey.Convey("/static/param route", func() {
				n := &node{}
				l := n.AddRoute(GET, "/foo/{name}", func(*Context) {})
				convey.So(n.find(GET, "/foo/bar"), convey.ShouldEqual, l)
			})
			convey.Convey("/static/param/static route", func() {
				n := &node{}
				l := n.AddRoute(GET, "/foo/{name}/bar", func(*Context) {})
				convey.So(n.find(GET, "/foo/jack/bar"), convey.ShouldEqual, l)
			})
			convey.Convey("static and param route", func() {
				n := &node{}
				l1 := n.AddRoute(GET, "/foo/{name}", func(*Context) {})
				l2 := n.AddRoute(GET, "/foo/bar", func(*Context) {})
				convey.So(n.find(GET, "/foo/bar"), convey.ShouldEqual, l2)
				convey.So(n.find(GET, "/foo/fuzz"), convey.ShouldEqual, l1)
			})
		})
		convey.Convey("Test Serve", func() {
			convey.Convey("Test All method", func() {
				n := &node{}
				r := RouterFunc(n.AddRoute)
				r.Options("/a/b", newSimpleHandler("OPTIONS"))
				r.Get("/a/b", newSimpleHandler("GET"))
				r.Head("/a/b", newSimpleHandler("HEAD"))
				r.Post("/a/b", newSimpleHandler("POST"))
				r.Put("/a/b", newSimpleHandler("PUT"))
				r.Delete("/a/b", newSimpleHandler("DELETE"))
				r.Trace("/a/b", newSimpleHandler("TRACE"))
				r.Connect("/a/b", newSimpleHandler("CONNECT"))
				r.Patch("/a/b", newSimpleHandler("PATCH"))
				r.Any("/a/c", newSimpleHandler("ANY"))
				convey.Convey("Test OPTIONS", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("OPTIONS", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "OPTIONS")
				})
				convey.Convey("Test Get", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "GET")
				})
				convey.Convey("Test HEAD", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("HEAD", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "HEAD")
				})
				convey.Convey("Test POST", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("POST", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "POST")
				})
				convey.Convey("Test PUT", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("PUT", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "PUT")
				})
				convey.Convey("Test DELETE", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("DELETE", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "DELETE")
				})
				convey.Convey("Test TRACE", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("TRACE", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "TRACE")
				})
				convey.Convey("Test CONNECT", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("CONNECT", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "CONNECT")
				})
				convey.Convey("Test PATCH", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("PATCH", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "PATCH")
				})
				convey.Convey("Test ANY", func() {
					convey.Convey("GET", func() {
						resp := httptest.NewRecorder()
						req, _ := http.NewRequest("GET", "/a/c", nil)
						ctx := newContext(defaultLogger)
						ctx.reset(NewResponseWriter(resp), req)
						n.Serve(ctx)
						convey.So(resp.Body.String(), convey.ShouldEqual, "ANY")
					})
					convey.Convey("POST", func() {
						resp := httptest.NewRecorder()
						req, _ := http.NewRequest("POST", "/a/c", nil)
						ctx := newContext(defaultLogger)
						ctx.reset(NewResponseWriter(resp), req)
						n.Serve(ctx)
						convey.So(resp.Body.String(), convey.ShouldEqual, "ANY")
					})
				})
			})

			convey.Convey("Test regexp path", func() {
				n := &node{}
				r := RouterFunc(n.AddRoute)
				r.Get(`/{name}`, newSimpleHandler("h1"))
				r.Get(`/{name}/{age}`, newSimpleHandler("h2"))
				convey.Convey("/jack", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/jack", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "h1")
				})
				convey.Convey("/jack/12", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/jack/12", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "h2")
				})
			})

			convey.Convey("Test error", func() {
				n := &node{}
				n.AddRoute(GET, "/a/b", newSimpleHandler("GET"))
				convey.Convey("test GET /a/b", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Body.String(), convey.ShouldEqual, "GET")
				})
				convey.Convey("test POST /a/b", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("POST", "/a/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Code, convey.ShouldEqual, http.StatusMethodNotAllowed)
				})
				convey.Convey("test GET /b", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/b", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
				})
				convey.Convey("test GET /a", func() {
					resp := httptest.NewRecorder()
					req, _ := http.NewRequest("GET", "/a", nil)
					ctx := newContext(defaultLogger)
					ctx.reset(NewResponseWriter(resp), req)
					n.Serve(ctx)
					convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
				})
			})
		})

		convey.Convey("Test param", func() {
			convey.Convey("/static/param/static route", func() {
				n := &node{}
				n.AddRoute(GET, "/hello/{name}", func(ctx *Context) {
					name, _ := ctx.Params().GetString("name")
					ctx.WriteString(name)
				})
				resp := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/hello/jack", nil)
				ctx := newContext(defaultLogger)
				ctx.reset(NewResponseWriter(resp), req)
				n.Serve(ctx)
				convey.So(resp.Body.String(), convey.ShouldEqual, "jack")
			})
		})
	})
}

func BenchmarkGinNode(b *testing.B) {
	n := &node{}
	n.AddRoute(GET, "/{a}/{b}/{c}/{d}/{e}", func(*Context) {})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		n.find(GET, "/a/b/c/d/e")
	}
}
