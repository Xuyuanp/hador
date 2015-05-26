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
	"compress/gzip"
	"net/http"
	"strings"
)

// AcceptEncoding is Accept-Encoding header
type AcceptEncoding string

// Accept checks if this coding could be accepted. See RFC2616-14.3
func (ae AcceptEncoding) Accept(coding string) bool {
	accept := false
	if len(ae) == 0 {
		if coding == "identity" ||
			coding == "gzip" ||
			coding == "compress" {
			return true
		}
		return false
	}
	codingslice := strings.Split(string(ae), ",")
	for _, codings := range codingslice {
		index := strings.Index(codings, ";")
		var co, q string
		if index != -1 {
			co = strings.TrimSpace(codings[:index])
			q = strings.TrimSpace(codings[index+1:][2:])
		} else {
			co = strings.TrimSpace(codings)
			q = ""
		}
		if co == "*" {
			if q == "0" {
				accept = false
			} else {
				accept = true
			}
		} else if co == coding {
			if q == "0" {
				accept = false
			} else {
				accept = true
			}
			break
		}
	}
	return accept
}

// GZipFilter returns a Filter that adds gzip compression to all requests.
// Make sure to use this before all other filters that alter the response body.
func GZipFilter(must bool) FilterFunc {
	return func(ctx *Context, next Handler) {
		acceptEncoding := AcceptEncoding(ctx.Request.Header.Get("Accept-Encoding"))
		if acceptEncoding.Accept("gzip") {
			ctx.Response.Header().Set("Content-Encoding", "gzip")
			grw := newGZipResponseWriter(ctx.Response)
			defer grw.Close()
			ctx.Response = grw
		} else if must {
			ctx.OnError(http.StatusNotAcceptable)
			return
		}

		next.Serve(ctx)
	}
}

type gzipResponseWriter struct {
	ResponseWriter
	w *gzip.Writer
}

func newGZipResponseWriter(rw ResponseWriter) *gzipResponseWriter {
	return &gzipResponseWriter{
		ResponseWriter: rw,
		w:              gzip.NewWriter(rw),
	}
}

func (grw *gzipResponseWriter) Write(data []byte) (int, error) {
	if "" == grw.Header().Get("Content-Type") {
		grw.Header().Set("Content-Type", http.DetectContentType(data))
	}
	size, err := grw.w.Write(data)
	return size, err
}

func (grw *gzipResponseWriter) Close() error {
	return grw.w.Close()
}
