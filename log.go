/*
 * Copyright 2015 Xuyuan Pang <xuyuanp # gmail dot com>
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
	"log"
	"os"
)

// Logger interface
type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})
	Critical(string, ...interface{})
}

type logger struct {
	*log.Logger
}

var defaultLogger = &logger{
	Logger: log.New(os.Stdout, "[Hodor] ", 0),
}

func (l *logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("\033[34m%s\033[0m", msg)
}

func (l *logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("\033[32m%s\033[0m", msg)
}

func (l *logger) Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("\033[36m%s\033[0m", msg)
}

func (l *logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("\033[33m%s\033[0m", msg)
}

func (l *logger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("\033[31m%s\033[0m", msg)
}
