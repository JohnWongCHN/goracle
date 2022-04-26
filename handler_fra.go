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

func fraHandler(ctx context.Context, conn OraClient, params map[string]string, _ ...string) (interface{}, error) {

	_sql := `
SELECT
	METRIC as "metric", 
	SUM(VALUE) AS "value"
FROM
	(
	SELECT
		'space_limit' AS METRIC, 
		SPACE_LIMIT AS VALUE
	FROM
		V$RECOVERY_FILE_DEST	
	UNION
	SELECT
		'space_used', 
		SPACE_USED AS VALUE
	FROM
		V$RECOVERY_FILE_DEST
	UNION
	SELECT
		'space_reclaimable', 
		SPACE_RECLAIMABLE AS VALUE
	FROM
		V$RECOVERY_FILE_DEST
	UNION
	SELECT
		'number_of_files', 
		NUMBER_OF_FILES AS VALUE
	FROM
		V$RECOVERY_FILE_DEST
	UNION
	SELECT
		'usable_pct', 
		DECODE(SPACE_LIMIT, 0, 0, (100 - (100 * (SPACE_USED - SPACE_RECLAIMABLE) / SPACE_LIMIT))) AS VALUE
	FROM
		V$RECOVERY_FILE_DEST
	UNION
	SELECT
		'restore_point', 
		COUNT(*) AS VALUE
	FROM
		V$RESTORE_POINT
	UNION
	SELECT
		DISTINCT *
	FROM
		TABLE(sys.ODCIVARCHAR2LIST('space_limit', 'space_used', 'space_reclaimable', 'number_of_files', 'usable_pct')), 
		TABLE(sys.ODCINUMBERLIST(0, 0, 0, 0, 0)) 
	)
GROUP BY
	METRIC
	`

	rows, err := conn.Query(ctx, _sql)
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
