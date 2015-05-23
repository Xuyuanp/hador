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
// If the handler's URL pattern is "/api/(?P<name>\w+)/{age:\d+}" and the request URL is
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

// GetUint gets params with key in uint format
func (params Params) GetUint(key string) (uint, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetUintMust gets params with key in uint format
func (params Params) GetUintMust(key string, def uint) uint {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return def
		}
		return uint(nu)
	}
	return def
}

// GetInt8 gets params with key in int8 format
func (params Params) GetInt8(key string) (int8, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetInt8Must gets params with key in int8 format
func (params Params) GetInt8Must(key string, def int8) int8 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return def
		}
		return int8(nu)
	}
	return def
}

// GetInt16 gets params with key in int16 format
func (params Params) GetInt16(key string) (int16, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return 0, err
		}
		return int16(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetInt16Must gets params with key in int16 format
func (params Params) GetInt16Must(key string, def int16) int16 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return def
		}
		return int16(nu)
	}
	return def
}

// GetInt32 gets params with key in int32 format
func (params Params) GetInt32(key string) (int32, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return 0, err
		}
		return int32(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetInt32Must gets params with key in int32 format
func (params Params) GetInt32Must(key string, def int32) int32 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(nu)
	}
	return def
}

// GetInt64 gets params with key in int64 format
func (params Params) GetInt64(key string) (int64, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		return int64(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetInt64Must gets params with key in int64 format
func (params Params) GetInt64Must(key string, def int64) int64 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return int64(nu)
	}
	return def
}

// GetUint8 gets params with key in uint8 format
func (params Params) GetUint8(key string) (uint8, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return 0, err
		}
		return uint8(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetUint8Must gets params with key in uint8 format
func (params Params) GetUint8Must(key string, def uint8) uint8 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return def
		}
		return uint8(nu)
	}
	return def
}

// GetUint16 gets params with key in uint16 format
func (params Params) GetUint16(key string) (uint16, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return 0, err
		}
		return uint16(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetUint16Must gets params with key in uint16 format
func (params Params) GetUint16Must(key string, def uint16) uint16 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return def
		}
		return uint16(nu)
	}
	return def
}

// GetUint32 gets params with key in uint32 format
func (params Params) GetUint32(key string) (uint32, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetUint32Must gets params with key in uint32 format
func (params Params) GetUint32Must(key string, def uint32) uint32 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return def
		}
		return uint32(nu)
	}
	return def
}

// GetUint64 gets params with key in uint64 format
func (params Params) GetUint64(key string) (uint64, error) {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint64(nu), nil
	}
	return 0, fmt.Errorf("No key named %s", key)
}

// GetUint64Must gets params with key in uint64 format
func (params Params) GetUint64Must(key string, def uint64) uint64 {
	if value, ok := params[key]; ok {
		nu, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return def
		}
		return uint64(nu)
	}
	return def
}

// GetBool gets params with key in boolean format
func (params Params) GetBool(key string) (bool, error) {
	if value, ok := params[key]; ok {
		b, err := strconv.ParseBool(value)
		return b, err
	}
	return false, fmt.Errorf("No key named %s", key)
}

// GetBoolMust gets params with key in boolean format
func (params Params) GetBoolMust(key string, def bool) bool {
	if value, ok := params[key]; ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return b
	}
	return def
}

// GetFloat32 gets params with key in float32 format
func (params Params) GetFloat32(key string) (float32, error) {
	if value, ok := params[key]; ok {
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0.0, err
		}
		return float32(f), err
	}
	return 0.0, fmt.Errorf("No key named %s", key)
}

// GetFloat32Must gets params with key in float32 format
func (params Params) GetFloat32Must(key string, def float32) float32 {
	if value, ok := params[key]; ok {
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(f)
	}
	return def
}

// GetFloat64 gets params with key in float64 format
func (params Params) GetFloat64(key string) (float64, error) {
	if value, ok := params[key]; ok {
		f, err := strconv.ParseFloat(value, 64)
		return f, err
	}
	return 0.0, fmt.Errorf("No key named %s", key)
}

// GetFloat64Must gets params with key in float64 format
func (params Params) GetFloat64Must(key string, def float64) float64 {
	if value, ok := params[key]; ok {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return f
	}
	return def
}
