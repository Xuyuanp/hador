/*
 * Copyright 2014 Xuyuan Pang <xuyuanp # gmail dot com>
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
	"time"
)

// NewLogFilter new a Filter to log each request details
func NewLogFilter(logger Logger) Filter {
	return FilterFunc(func(ctx *Context, next Handler) {
		start := time.Now()
		req := ctx.Request
		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}
		path := req.URL.Path
		if req.URL.RawQuery != "" {
			path += req.URL.RawQuery
		}
		logger.Info("Started %s %s %s", req.Method, path, addr)

		next.Serve(ctx)

		rw := ctx.Response
		status := rw.Status()
		statusText := http.StatusText(status)
		duration := time.Since(start)
		if status >= 500 {
			logger.Critical("Completed %d %s in %v", status, statusText, duration)
		} else if status >= 400 {
			logger.Error("Completed %d %s in %v", status, statusText, duration)
		} else if status >= 300 {
			logger.Warning("Completed %d %s in %v", status, statusText, duration)
		} else if status >= 200 {
			logger.Info("Completed %d %s in %v", status, statusText, duration)
		} else {
			logger.Debug("Completed %d %s in %v", status, statusText, duration)
		}
	})
}
