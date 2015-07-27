package hador

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNode(t *testing.T) {
	convey.Convey("Test node", t, func() {
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

		convey.Convey("Test match", func() {
			convey.Convey("single", func() {
				n := &node{}
				l := n.AddRoute(GET, "/bar", func(_ *Context) {})
				_, lr, err := n.match(GET, "/bar", nil)
				convey.So(err, convey.ShouldBeNil)
				convey.So(lr, convey.ShouldEqual, l)
			})
			convey.Convey("static", func() {
				n := &node{}
				l := n.AddRoute(GET, "/foo/bar", func(_ *Context) {})
				_, lr, err := n.match(GET, "/foo/bar", nil)
				convey.So(err, convey.ShouldBeNil)
				convey.So(lr, convey.ShouldEqual, l)
			})
			convey.Convey("param", func() {
				n := &node{}
				l := n.AddRoute(GET, "/hello/{name}", func(_ *Context) {})
				maxParams := n.findMaxParams()
				convey.So(maxParams, convey.ShouldEqual, 1)
				params := make(Params, maxParams)
				_, lr, err := n.match(GET, "/hello/jack", params[0:0])
				convey.So(err, convey.ShouldBeNil)
				convey.So(lr, convey.ShouldEqual, l)
			})
			convey.Convey("param with regexp", func() {
				n := &node{}
				l := n.AddRoute(GET, "/hello/{name:\\d+}", func(_ *Context) {})
				maxParams := n.findMaxParams()
				convey.So(maxParams, convey.ShouldEqual, 1)
				params := make(Params, maxParams)
				_, lr, err := n.match(GET, "/hello/jack", params[0:0])
				convey.So(err, convey.ShouldEqual, err404)

				_, lr, err = n.match(GET, "/hello/123", params[0:0])
				convey.So(err, convey.ShouldBeNil)
				convey.So(lr, convey.ShouldEqual, l)
			})
		})

		convey.Convey("Test findMaxParams", func() {
			n := &node{}
			n.AddRoute(GET, "/foo/bar", func(_ *Context) {})
			n.AddRoute(GET, "/foo/{b}", func(_ *Context) {})
			n.AddRoute(GET, "/hello/{a}/{b}", func(_ *Context) {})
			n.AddRoute(GET, "/{a}/{b}/{c}", func(_ *Context) {})
			convey.So(n.findMaxParams(), convey.ShouldEqual, 3)
		})

		convey.Convey("Test readField", func() {
			convey.Convey("Test end", func() {
				str := `{age:\d+:integer:your age}`
				field, rest, end := readField(str[1:])
				convey.So(field, convey.ShouldEqual, "age")
				convey.So(rest, convey.ShouldEqual, `\d+:integer:your age}`)
				convey.So(end, convey.ShouldBeFalse)

				field, rest, end = readField(rest)
				convey.So(field, convey.ShouldEqual, `\d+`)
				convey.So(rest, convey.ShouldEqual, `integer:your age}`)
				convey.So(end, convey.ShouldBeFalse)

				field, rest, end = readField(rest)
				convey.So(field, convey.ShouldEqual, `integer`)
				convey.So(rest, convey.ShouldEqual, `your age}`)
				convey.So(end, convey.ShouldBeFalse)

				field, rest, end = readField(rest)
				convey.So(field, convey.ShouldEqual, `your age`)
				convey.So(rest, convey.ShouldEqual, ``)
				convey.So(end, convey.ShouldBeTrue)
			})
			convey.Convey("Test not end", func() {
				str := `{age:\d+:integer:your age}/hello`
				field, rest, end := readField(str[1:])
				convey.So(field, convey.ShouldEqual, "age")
				convey.So(rest, convey.ShouldEqual, `\d+:integer:your age}/hello`)
				convey.So(end, convey.ShouldBeFalse)

				field, rest, end = readField(rest)
				convey.So(field, convey.ShouldEqual, `\d+`)
				convey.So(rest, convey.ShouldEqual, `integer:your age}/hello`)
				convey.So(end, convey.ShouldBeFalse)

				field, rest, end = readField(rest)
				convey.So(field, convey.ShouldEqual, `integer`)
				convey.So(rest, convey.ShouldEqual, `your age}/hello`)
				convey.So(end, convey.ShouldBeFalse)

				field, rest, end = readField(rest)
				convey.So(field, convey.ShouldEqual, `your age`)
				convey.So(rest, convey.ShouldEqual, `/hello`)
				convey.So(end, convey.ShouldBeTrue)
			})
		})

		convey.Convey("Test readParam", func() {
			str := `{age:\d+:integer:your age}/hello`
			name, regstr, dataType, desc, rest := readParam(str)
			convey.So(name, convey.ShouldEqual, "age")
			convey.So(regstr, convey.ShouldEqual, `\d+`)
			convey.So(dataType, convey.ShouldEqual, "integer")
			convey.So(desc, convey.ShouldEqual, "your age")
			convey.So(rest, convey.ShouldEqual, "/hello")
		})
	})
}
