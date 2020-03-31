// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmt

import "errors"

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
//
// If the format specifier includes a %w verb with an error operand,					// 如果格式中包括%w动词，并且存在相应的error参数，返回的error将会实现Unwrap方法且该方法会返回原error参数。
// the returned error will implement an Unwrap method returning the operand. It is
// invalid to include more than one %w verb or to supply it with an operand				// 最多支持1个%w动词，且%w对应的参数必须实现error接口。否则%w则将退化成%v。
// that does not implement the error interface. The %w verb is otherwise
// a synonym for %v.
func Errorf(format string, a ...interface{}) error {
	p := newPrinter()
	p.wrapErrs = true
	p.doPrintf(format, a) // 解析格式，如果发现%w动词且提供了合法的error参数，则把error放到p.wrappedErr成员中
	s := string(p.buf)
	var err error
	if p.wrappedErr == nil { // 没有%w动词，生成基础error
		err = errors.New(s)
	} else { // 存在%w动词，生成wrapError
		err = &wrapError{s, p.wrappedErr}
	}
	p.free()
	return err
}

type wrapError struct {
	msg string // 存储上下文信息和err.Error()
	err error  // 存储原error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}
