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
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	headerContentType     = "Content-Type"
	headerContentEncoding = "Content-Encoding"
	headerAcceptEncoding  = "Accept-Encoding"

	encodingTypeGZip     = "gzip"
	encodingTypeDeflate  = "deflate"
	encodingTypeIdentity = "identity"
	encodingTypeCompress = "compress"
)

// AcceptEncoding is Accept-Encoding header
type AcceptEncoding string

// Accept checks if this coding could be accepted. See RFC2616-14.3
func (ae AcceptEncoding) Accept(coding string) bool {
	accept := false
	if len(ae) == 0 {
		if coding == encodingTypeIdentity ||
			coding == encodingTypeGZip ||
			coding == encodingTypeCompress {
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

type encodingWriter struct {
	ResponseWriter
	w io.WriteCloser
}

func newEncodingWriter(rw ResponseWriter, wrapper func(io.Writer) io.WriteCloser) *encodingWriter {
	return &encodingWriter{
		ResponseWriter: rw,
		w:              wrapper(rw),
	}
}

func (ew *encodingWriter) Write(p []byte) (int, error) {
	if "" == ew.Header().Get(headerContentType) {
		ew.Header().Set(headerContentType, http.DetectContentType(p))
	}
	size, err := ew.w.Write(p)
	return size, err
}

func (ew *encodingWriter) Close() error {
	return ew.w.Close()
}

// GZipFilter returns a Filter that adds gzip compression to all requests.
// Make sure to use this before all other filters that alter the response body.
func GZipFilter(must bool) FilterFunc {
	return func(ctx *Context, next Handler) {
		acceptEncoding := AcceptEncoding(ctx.Request.Header.Get(headerAcceptEncoding))
		if acceptEncoding.Accept(encodingTypeGZip) {
			ctx.Response.Header().Set(headerContentEncoding, encodingTypeGZip)
			ew := newEncodingWriter(ctx.Response, func(w io.Writer) io.WriteCloser {
				return gzip.NewWriter(w)
			})
			defer ew.Close()
			ctx.Response = ew
		} else if must {
			ctx.OnError(http.StatusNotAcceptable)
			return
		}

		next.Serve(ctx)
	}
}

// DeflateFilter returns a Filter that adds flate compression to all requests.
// Make sure to use this before all other filters that alter the response body.
func DeflateFilter(level int, must bool) FilterFunc {
	return func(ctx *Context, next Handler) {
		acceptEncoding := AcceptEncoding(ctx.Request.Header.Get(headerAcceptEncoding))
		if acceptEncoding.Accept(encodingTypeDeflate) {
			ctx.Response.Header().Set(headerContentEncoding, encodingTypeDeflate)
			ew := newEncodingWriter(ctx.Response, func(w io.Writer) io.WriteCloser {
				fw, _ := flate.NewWriter(w, level)
				return fw
			})
			defer ew.Close()
			ctx.Response = ew
		} else if must {
			ctx.OnError(http.StatusNotAcceptable)
			return
		}

		next.Serve(ctx)
	}
}
