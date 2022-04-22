/*
** Zabbix
** Copyright (C) 2001-2022 Zabbix SIA
**
** This program is free software; you can redistribute it and/or modify
** it under the terms of the GNU General Public License as published by
** the Free Software Foundation; either version 2 of the License, or
** (at your option) any later version.
**
** This program is distributed in the hope that it will be useful,
** but WITHOUT ANY WARRANTY; without even the implied warranty of
** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
** GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License
** along with this program; if not, write to the Free Software
** Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
**/

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
