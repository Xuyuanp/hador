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
	"strconv"
)

// Params is a wrapper of map[string]string to handle the params in regexp pattern.
// If the handler's URL pattern is "/api/(?P<name>\w+)/(?P<age>\d+)" and the request URL is
// "/api/jack/12", the Params will be { "name": "jack", "age": "12" }.
// Using `hodor.GetParams(req)` to get the Params.
type Params map[string]string

// GetString method returns the param named `key` as string type.
func (params Params) GetString(key string) (string, error) {
	if value, ok := params[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("No key named %s", key)
}

// GetStringMust method returns the param named `key` as string type.
// If no param with the provided name, the default value will be returned.
func (params Params) GetStringMust(key string, def string) string {
	if value, ok := params[key]; ok {
		return value
	}
	return def
}

// GetInt method does the same work as GetString, but converts the string into integer.
// If the param is not a valid integer, an err will be returned.
func (params Params) GetInt(key string) (int, error) {
	if value, ok := params[key]; ok {
		return strconv.Atoi(value)
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetIntMust method does the same work as GetStringMust, but converts the string into integer.
// If the param is not a valid integer, the default value will be returned.
func (params Params) GetIntMust(key string, def int) int {
	if value, ok := params[key]; ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return def
}
