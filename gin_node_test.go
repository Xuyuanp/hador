package hador

import (
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
		convey.Convey("Test param", func() {
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
	})
}
