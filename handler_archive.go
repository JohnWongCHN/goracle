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

func archiveHandler(ctx context.Context, conn OraClient, params map[string]string, _ ...string) (interface{}, error) {
	_sql := `
SELECT
	d.DEST_NAME as "dest_name",
	DECODE(d.STATUS, 'VALID', 3, 'DEFERRED', 2, 'ERROR', 1, 0) as "status",
	d.LOG_SEQUENCE as "log_sequence",
	NVL(TO_CHAR(d.ERROR), ' ')	as "error"
FROM
	V$ARCHIVE_DEST d,
	V$DATABASE db
WHERE
	d.STATUS != 'INACTIVE'
	AND db.LOG_MODE = 'ARCHIVELOG'
	`

	rows, err := conn.Query(ctx, _sql)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}
	defer rows.Close()

	// JSON marshaling
	var data []string

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

		for i, value := range values {
			// skip dest_name column
			if i == 0 {
				continue
			}
			results[columns[i]] = value
		}

		// generate proper map
		_data := map[string]interface{}{
			fmt.Sprintf("%v", values[0]): results}

		// jsonRes, _ := json.Marshal(results)
		jsonRes, _ := json.Marshal(_data)
		data = append(data, strings.TrimSpace(string(jsonRes)))
	}

	return "[" + strings.Join(data, ",") + "]", nil
}
