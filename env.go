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
	"os"
)

const (
	// DEVELOP environment
	DEVELOP = "develop"

	// PRODUCTION environment
	PRODUCTION = "production"
)

// ENV is runtime environment, develop by default
var ENV = DEVELOP

func init() {
	env := os.Getenv("HADOR_ENV")
	if env != "" {
		if env != DEVELOP && env != PRODUCTION {
			panic(fmt.Errorf("Invalid env: %s, choose %s or %s.", env, DEVELOP, PRODUCTION))
		}
		ENV = env
	}
}
