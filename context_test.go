/*
 * Copyright 2015 Xuyuan Pang
 * Author: Xuyuan Pang
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

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {
	convey.Convey("given a new context", t, func() {
		ctx := newContext(defaultLogger)
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		convey.So(err, convey.ShouldBeNil)
		rw := NewResponseWriter(resp)
		ctx.reset(rw, req)

		convey.Convey("test internal data", func() {
			convey.Convey("get & delete should be nil", func() {
				convey.So(ctx.Get("foo"), convey.ShouldBeNil)
				_, ok := ctx.GetOK("foo")
				convey.So(ok, convey.ShouldBeFalse)
				convey.So(ctx.Delete("foo"), convey.ShouldBeNil)
				convey.Convey("set key-value", func() {
					ctx.Set("foo", "bar")
					convey.Convey("get & delete should be ok", func() {
						convey.So(ctx.Get("foo"), convey.ShouldEqual, "bar")
						_, ok := ctx.GetOK("foo")
						convey.So(ok, convey.ShouldBeTrue)

						convey.So(ctx.Delete("foo"), convey.ShouldEqual, "bar")
					})
					convey.Convey("get & delete unknown key should be nil", func() {
						convey.So(ctx.Get("bazz"), convey.ShouldBeNil)
						_, ok := ctx.GetOK("bazz")
						convey.So(ok, convey.ShouldBeFalse)

						convey.So(ctx.Delete("bazz"), convey.ShouldBeNil)
					})
				})
			})
		})

		convey.Convey("test error handler", func() {
			convey.Convey("test custome error handler", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)
				ctx.SetErrorHandler(http.StatusNotFound, func(...interface{}) {
					ctx.WriteHeader(http.StatusNotFound)
					ctx.Response.Write([]byte("404"))
				})
				ctx.OnError(http.StatusNotFound)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
				convey.So(resp.Body.String(), convey.ShouldEqual, "404")
			})
			convey.Convey("test default error handler", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)
				ctx.OnError(http.StatusNotFound)
				convey.So(resp.Code, convey.ShouldEqual, http.StatusNotFound)
				convey.So(resp.Body.String(), convey.ShouldEqual,
					http.StatusText(http.StatusNotFound)+"\n")
			})
		})

		convey.Convey("test WriteString", func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			convey.So(err, convey.ShouldBeNil)
			rw := NewResponseWriter(resp)
			ctx.reset(rw, req)

			content := "foobar"
			_, err = ctx.WriteString(content)
			convey.So(err, convey.ShouldBeNil)
			convey.So(resp.Body.String(), convey.ShouldEqual, content)
		})

		convey.Convey("test render", func() {
			type cv struct {
				Foo string
				Bar int
			}
			v := cv{
				Foo: "blah",
				Bar: 10,
			}
			jsondata, err := json.Marshal(v)
			convey.So(err, convey.ShouldBeNil)
			xmldata, err := xml.Marshal(v)
			convey.So(err, convey.ShouldBeNil)

			jsonindentdata, err := json.MarshalIndent(v, "", "\t")
			convey.So(err, convey.ShouldBeNil)
			xmlindentdata, err := xml.MarshalIndent(v, "", "\t")
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("test RenderJSON", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)

				err = ctx.RenderJSON(v, 201)
				convey.So(err, convey.ShouldBeNil)
				convey.So(resp.Code, convey.ShouldEqual, 201)
				convey.So(resp.Header().Get("Content-Type"), convey.ShouldEqual,
					"application/json; charset=utf-8")
				convey.So(resp.Body.String(), convey.ShouldEqual, string(jsondata)+"\n")
			})
			convey.Convey("test RenderPrettyJSON", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)

				err = ctx.RenderPrettyJSON(v)
				convey.So(err, convey.ShouldBeNil)
				convey.So(resp.Code, convey.ShouldEqual, 200)
				convey.So(resp.Header().Get("Content-Type"), convey.ShouldEqual,
					"application/json; charset=utf-8")
				convey.So(resp.Body.String(), convey.ShouldEqual, string(jsonindentdata))
			})
			convey.Convey("test RenderXML", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)

				err = ctx.RenderXML(v, 201)
				convey.So(err, convey.ShouldBeNil)
				convey.So(resp.Code, convey.ShouldEqual, 201)
				convey.So(resp.Header().Get("Content-Type"), convey.ShouldEqual,
					"application/xml; charset=utf-8")
				convey.So(resp.Body.String(), convey.ShouldEqual,
					string(xmldata))

				convey.Convey("test RenderXML error", func() {
					resp := httptest.NewRecorder()
					req, err := http.NewRequest("GET", "/", nil)
					convey.So(err, convey.ShouldBeNil)
					rw := NewResponseWriter(resp)
					ctx.reset(rw, req)

					v := map[string]string{}

					err = ctx.RenderXML(v, 201)
					convey.So(err, convey.ShouldNotBeNil)
				})
			})
			convey.Convey("test RenderPrettyXML", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/", nil)
				convey.So(err, convey.ShouldBeNil)
				rw := NewResponseWriter(resp)
				ctx.reset(rw, req)

				err = ctx.RenderPrettyXML(v)
				convey.So(err, convey.ShouldBeNil)
				convey.So(resp.Code, convey.ShouldEqual, 200)
				convey.So(resp.Header().Get("Content-Type"), convey.ShouldEqual,
					"application/xml; charset=utf-8")
				convey.So(resp.Body.String(), convey.ShouldEqual, string(xmlindentdata))
			})
			convey.Convey("test resolve methods", func() {
				convey.Convey("test resolve json", func() {
					jsonData, _ := json.Marshal(v)
					r := bytes.NewReader(jsonData)
					resp := httptest.NewRecorder()
					req, err := http.NewRequest("POST", "/", r)
					convey.So(err, convey.ShouldBeNil)
					rw := NewResponseWriter(resp)
					ctx.reset(rw, req)

					j := cv{}
					err = ctx.ResolveJSON(&j)
					convey.So(err, convey.ShouldBeNil)
					convey.So(j.Bar, convey.ShouldEqual, v.Bar)
					convey.So(j.Foo, convey.ShouldEqual, v.Foo)
				})
				convey.Convey("test resolve xml", func() {
					xmlData, _ := xml.Marshal(v)
					r := bytes.NewReader(xmlData)
					resp := httptest.NewRecorder()
					req, err := http.NewRequest("POST", "/", r)
					convey.So(err, convey.ShouldBeNil)
					rw := NewResponseWriter(resp)
					ctx.reset(rw, req)

					x := cv{}
					err = ctx.ResolveXML(&x)
					convey.So(err, convey.ShouldBeNil)
					convey.So(x.Bar, convey.ShouldEqual, v.Bar)
					convey.So(x.Foo, convey.ShouldEqual, v.Foo)
				})
			})

		})

	})
}
