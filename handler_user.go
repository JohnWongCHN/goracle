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
	"encoding/json"

	"git.zabbix.com/ap/plugin-support/zbxerr"
)

func userHandler(ctx context.Context, conn OraClient, params map[string]string, _ ...string) (interface{}, error) {
	var res int

	username := conn.WhoAmI()
	if params["Username"] != "" {
		username = params["Username"]
	}

	_sql := `
SELECT
	ROUND(DECODE(SIGN(NVL(EXPIRY_DATE, SYSDATE + 999) - SYSDATE), -1, 0, NVL(EXPIRY_DATE, SYSDATE + 999) - SYSDATE)) as "exp_passwd_days_before"
FROM
	DBA_USERS
WHERE
	USERNAME = UPPER(:1)
	`

	row, err := conn.QueryRow(ctx, _sql, username)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	err = row.Scan(&res)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	jsonRes, _ := json.Marshal(map[string]int{"exp_passwd_days_before": res})

	return string(jsonRes), nil
}
