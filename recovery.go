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
	"fmt"
	"net/http"
	"runtime"
)

// NewRecoveryFilter return a Filter to recover all unrecovered panic
func NewRecoveryFilter(logger Logger) FilterFunc {
	return func(ctx *Context, next Handler) {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				stack := trace[:n]
				msg := fmt.Sprintf("PANIC: %v\n%s", err, stack)
				logger.Critical(msg)
				ctx.OnError(http.StatusInternalServerError, err)
			}
		}()

		next.Serve(ctx)
	}
}
