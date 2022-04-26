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
	"context"
	"fmt"
)

const (
	pingFailed = 0
	pingOk     = 1
)

// pingHandler queries 'SELECT 1 FROM DUAL' and returns pingOk if a connection is alive or pingFailed otherwise.
func pingHandler(ctx context.Context, conn OraClient, params map[string]string, _ ...string) (interface{}, error) {
	var res int

	row, err := conn.QueryRow(ctx, fmt.Sprintf("SELECT %d FROM DUAL", pingOk))
	if err != nil {
		return pingFailed, nil
	}

	err = row.Scan(&res)

	if err != nil || res != pingOk {
		return pingFailed, nil
	}

	return pingOk, nil
}
