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
		})
		convey.Convey("Test find", func() {
			convey.Convey("Test static", func() {
				convey.Convey("root", func() {
					n := &node{}
					n.AddRoute(GET, "/", func(*Context) {})
					convey.So(n.find(GET, "/"), convey.ShouldNotBeNil)
				})
				convey.Convey("single route", func() {
					n := &node{}
					n.AddRoute(GET, "/foo", func(*Context) {})
					convey.So(n.find(GET, "/foo"), convey.ShouldNotBeNil)
				})
				convey.Convey("multi route with same prefix", func() {
					n := &node{}
					n.AddRoute(GET, "/initfoo", func(*Context) {})
					n.AddRoute(GET, "/initbar", func(*Context) {})
					convey.So(n.find(GET, "/initfoo"), convey.ShouldNotBeNil)
					convey.So(n.find(GET, "/initbar"), convey.ShouldNotBeNil)
				})
				convey.Convey("multi routes without same prefix", func() {
					n := &node{}
					n.AddRoute(GET, "/foo", func(*Context) {})
					n.AddRoute(GET, "/bar", func(*Context) {})
					convey.So(n.find(GET, "/foo"), convey.ShouldNotBeNil)
				})
				convey.Convey("multi long routes", func() {
					n := &node{}
					n.AddRoute(GET, "/abc/def/ghi/jkl/mno/pqr", func(*Context) {})
					n.AddRoute(GET, "/abc/xyz/ghi/jkl/mno/pqr", func(*Context) {})
					convey.So(n.find(GET, "/abc/xyz/ghi/jkl/mno/pqr"), convey.ShouldNotBeNil)
				})
			})
		})
	})
}
