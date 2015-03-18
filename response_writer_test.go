/*
 * Copyright (c) 2014 Jeremy Saenz
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package hador

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))

	convey.Convey("TestResponseWriterWritingString", t, func() {
		convey.So(rec.Code, convey.ShouldEqual, rw.Status())
		convey.So(rec.Body.String(), convey.ShouldEqual, "Hello world")
		convey.So(rw.Status(), convey.ShouldEqual, http.StatusOK)
		convey.So(rw.Size(), convey.ShouldEqual, 11)
		convey.So(rw.Written(), convey.ShouldEqual, true)
	})
}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))
	rw.Write([]byte("foo bar bat baz"))
	convey.Convey("TestResponseWriterWritingStrings", t, func() {
		convey.So(rec.Code, convey.ShouldEqual, rw.Status())
		convey.So(rec.Body.String(), convey.ShouldEqual, "Hello worldfoo bar bat baz")
		convey.So(rw.Status(), convey.ShouldEqual, http.StatusOK)
		convey.So(rw.Size(), convey.ShouldEqual, 26)
	})
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)

	convey.Convey("TestResponseWriterWritingHeader", t, func() {
		convey.So(rec.Code, convey.ShouldEqual, rw.Status())
		convey.So(rec.Body.String(), convey.ShouldEqual, "")
		convey.So(rw.Status(), convey.ShouldEqual, http.StatusNotFound)
		convey.So(rw.Size(), convey.ShouldEqual, 0)
	})
}

func TestResponseWriterBefore(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)
	result := ""

	rw.Before(func(ResponseWriter) {
		result += "foo"
	})
	rw.Before(func(ResponseWriter) {
		result += "bar"
	})

	rw.WriteHeader(http.StatusNotFound)

	convey.Convey("TestResponseWriterBefore", t, func() {
		convey.So(rec.Code, convey.ShouldEqual, rw.Status())
		convey.So(rec.Body.String(), convey.ShouldEqual, "")
		convey.So(rw.Status(), convey.ShouldEqual, http.StatusNotFound)
		convey.So(rw.Size(), convey.ShouldEqual, 0)
		convey.So(result, convey.ShouldEqual, "barfoo")
	})
}

func TestResponseWriterHijack(t *testing.T) {
	convey.Convey("TestResponseWriterHijack", t, func() {
		hijackable := newHijackableResponse()
		rw := NewResponseWriter(hijackable)
		hijacker, ok := rw.(http.Hijacker)

		convey.So(ok, convey.ShouldBeTrue)
		_, _, err := hijacker.Hijack()
		if err != nil {
			t.Error(err)
		}
		convey.So(hijackable.Hijacked, convey.ShouldBeTrue)
	})
}

func TestResponseWriteHijackNotOK(t *testing.T) {
	convey.Convey("TestResponseWriteHijackNotOK", t, func() {
		hijackable := new(http.ResponseWriter)
		rw := NewResponseWriter(*hijackable)
		hijacker, ok := rw.(http.Hijacker)
		convey.So(ok, convey.ShouldBeTrue)
		_, _, err := hijacker.Hijack()

		convey.So(err, convey.ShouldNotBeNil)
	})
}

func TestResponseWriterCloseNotify(t *testing.T) {
	convey.Convey("TestResponseWriterCloseNotify", t, func() {
		rec := newCloseNotifyingRecorder()
		rw := NewResponseWriter(rec)
		closed := false
		notifier := rw.(http.CloseNotifier).CloseNotify()
		rec.close()
		select {
		case <-notifier:
			closed = true
		case <-time.After(time.Second):
		}
		convey.So(closed, convey.ShouldBeTrue)
	})
}

func TestResponseWriterFlusher(t *testing.T) {
	convey.Convey("TestResponseWriterFlusher", t, func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)

		f, ok := rw.(http.Flusher)
		convey.So(ok, convey.ShouldBeTrue)
		f.Flush()
	})
}
