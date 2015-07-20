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
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLeaf(t *testing.T) {
	convey.Convey("Test Leaf", t, func() {
		h := New()
		handler := newSimpleHandler("swagger")
		leaf := h.Get("/swagger", handler)
		convey.So(leaf, convey.ShouldNotBeNil)
		convey.So(leaf.Handler(), convey.ShouldEqual, handler)
		convey.So(leaf.Method(), convey.ShouldEqual, GET)
		parent := leaf.parent
		for parent != nil && parent.parent != nil {
			parent = parent.parent
		}
		convey.So(parent, convey.ShouldEqual, h.root)
	})
}
