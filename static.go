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
	"net/http"
	"strings"
)

// Static struct
type Static struct {
	Prefix    string
	IndexFile string
	Dir       http.FileSystem
}

// NewStatic creates new Static instance
func NewStatic(dir http.FileSystem) *Static {
	return &Static{
		Prefix:    "",
		IndexFile: "index.html",
		Dir:       dir,
	}
}

// Filter serves all static file request
func (s *Static) Filter(ctx *Context, next Handler) {
	if ctx.Request.Method != "GET" && ctx.Request.Method != "HEAD" {
		next.Serve(ctx)
		return
	}
	path := ctx.Request.URL.Path
	if s.Prefix != "" {
		if !strings.HasPrefix(path, s.Prefix) {
			next.Serve(ctx)
			return
		}
		path = path[len(s.Prefix):]
		if path != "" && path[0] != '/' {
			next.Serve(ctx)
			return
		}
	}

	file, err := s.Dir.Open(path)
	if err != nil {
		next.Serve(ctx)
		return
	}
	defer file.Close()

	fs, err := file.Stat()
	if err != nil {
		next.Serve(ctx)
		return
	}

	// serve index file
	if fs.IsDir() {
		if !strings.HasSuffix(ctx.Request.URL.Path, "/") {
			http.Redirect(ctx.Response, ctx.Request.Request, ctx.Request.URL.Path+"/", http.StatusFound)
			return
		}

		path = path + s.IndexFile
		file, err := s.Dir.Open(path)
		if err != nil {
			next.Serve(ctx)
			return
		}
		defer file.Close()

		fs, err = file.Stat()
		if err != nil || fs.IsDir() {
			next.Serve(ctx)
			return
		}
	}

	http.ServeContent(ctx.Response, ctx.Request.Request, path, fs.ModTime(), file)
}
