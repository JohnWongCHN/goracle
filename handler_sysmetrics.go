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
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"git.zabbix.com/ap/plugin-support/zbxerr"
)

const (
	duration60sec = "2"
	duration15sec = "3"
)

func sysMetricsHandler(ctx context.Context, conn OraClient, params map[string]string,
	_ ...string) (interface{}, error) {

	var (
		groupID = duration60sec
	)

	switch params["Duration"] {
	case "15":
		groupID = duration15sec
	case "60":
		groupID = duration60sec
	default:
		return nil, zbxerr.ErrorInvalidParams
	}

	_sql := `
SELECT
	METRIC_NAME as "metric_name",
	ROUND(VALUE, 3) as "value"
FROM
	V$SYSMETRIC
WHERE
	GROUP_ID = :1
	`

	rows, err := conn.Query(ctx, _sql, groupID)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(values))

	for i := range values {
		valuePointers[i] = &values[i]
	}

	results := make(map[string]interface{})

	for rows.Next() {
		err = rows.Scan(valuePointers...)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, zbxerr.ErrorEmptyResult.Wrap(err)
			}

			return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
		}

		results[fmt.Sprintf("%v", values[0])] = values[1]

	}
	jsonRes, _ := json.Marshal(results)

	return strings.TrimSpace(string(jsonRes)), nil
}
