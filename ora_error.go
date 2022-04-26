/*
MIT License

Copyright (c) [2022] [John Wong<john-wong@outlook.com>]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
)

// OraErr is an error holding the ORA-01234 code and the message.
type OraErr struct {
	message, funName, action, sqlState string
	code, offset                       int
	recoverable, warning               bool
}

// AsOraErr returns the underlying *OraErr and whether it succeeded.
func AsOraErr(err error) (*OraErr, bool) {
	var oerr *OraErr
	ok := errors.As(err, &oerr)
	return oerr, ok
}

var _ error = (*OraErr)(nil)

// Code returns the OraErr's error code.
func (oe *OraErr) Code() int {
	if oe == nil {
		return 0
	}
	return oe.code
}

// Message returns the OraErr's message.
func (oe *OraErr) Message() string {
	if oe == nil {
		return ""
	}
	return oe.message
}
func (oe *OraErr) Error() string {
	if oe == nil {
		return ""
	}
	msg := oe.Message()
	if oe.code == 0 && msg == "" {
		return ""
	}
	return fmt.Sprintf("ORA-%05d: %s", oe.code, oe.message)
}

// Offset returns the parse error offset (in bytes) when executing a statement or the row offset when performing bulk operations or fetching batch error information. If neither of these cases are true, the value is 0.
func (oe *OraErr) Offset() int { return oe.offset }

// FunName returns the public ODPI-C function name which was called in which the error took place. This is a null-terminated ASCII string.
func (oe *OraErr) FunName() string { return oe.funName }

// Action returns the internal action that was being performed when the error took place. This is a null-terminated ASCII string.
func (oe *OraErr) Action() string { return oe.action }

// SQLState returns the SQLSTATE code associated with the error. This is a 5 character null-terminated string.
func (oe *OraErr) SQLState() string { return oe.sqlState }

// Recoverable indicates if the error is recoverable. This is always false unless both client and server are at release 12.1 or higher.
func (oe *OraErr) Recoverable() bool { return oe.recoverable }

func (oe *OraErr) IsWarning() bool { return oe.warning }
