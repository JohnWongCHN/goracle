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

func sgaHandler(ctx context.Context, conn OraClient, params map[string]string, _ ...string) (interface{}, error) {

	_sql := `
SELECT
	POOL as "pool", 
	SUM(BYTES) AS "bytes"
FROM
	(
	SELECT
		LOWER(REPLACE(POOL, ' ', '_')) AS POOL,
		SUM(BYTES) AS BYTES
	FROM
		V$SGASTAT
	WHERE
		POOL IN ('java pool', 'large pool')
	GROUP BY
		POOL
	UNION
	SELECT
		'shared_pool',
		SUM(BYTES)
	FROM
		V$SGASTAT
	WHERE
		POOL = 'shared pool'
		AND NAME NOT IN ('library cache', 'dictionary cache', 'free memory', 'sql area')
	UNION
	SELECT
		NAME,
		BYTES
	FROM
		V$SGASTAT
	WHERE
		POOL IS NULL
		AND NAME IN ('log_buffer', 'fixed_sga')
	UNION
	SELECT
		'buffer_cache',
		SUM(BYTES)
	FROM
		V$SGASTAT
	WHERE
		POOL IS NULL
		AND NAME IN ('buffer_cache', 'db_block_buffers')
	UNION
	SELECT
		DISTINCT *
	FROM
		TABLE(sys.ODCIVARCHAR2LIST('buffer_cache', 'fixed_sga', 'java_pool', 'large_pool', 'log_buffer', 'shared_pool')), 
		TABLE(sys.ODCINUMBERLIST(0, 0, 0, 0, 0, 0))	
	)
GROUP BY
	POOL
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
