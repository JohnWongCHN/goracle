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

func sessionsHandler(ctx context.Context, conn OraClient, params map[string]string, _ ...string) (interface{}, error) {

	_sql := `
SELECT
	METRIC, SUM(VALUE) AS VALUE
FROM
	(
	SELECT
		LOWER(REPLACE(STATUS || ' ' || TYPE, ' ', '_')) AS METRIC, 
		COUNT(*) AS VALUE
	FROM
		V$SESSION
	GROUP BY
		STATUS, TYPE
	UNION
	SELECT
		DISTINCT *
	FROM
		TABLE(sys.ODCIVARCHAR2LIST('inactive_user', 'active_user', 'active_background')), 
		TABLE(sys.ODCINUMBERLIST(0, 0, 0))
	)
GROUP BY
	METRIC
UNION
SELECT
	'total' AS METRIC, 
	COUNT(*) AS VALUE
FROM
	V$SESSION 
UNION
SELECT
	'long_time_locked' AS METRIC, 
	COUNT(*) AS VALUE
FROM
	V$SESSION
WHERE
	BLOCKING_SESSION IS NOT NULL
	AND BLOCKING_SESSION_STATUS = 'VALID'
	AND SECONDS_IN_WAIT > :1
UNION
SELECT
	'lock_rate' ,
	(CNT_BLOCK / CNT_ALL) * 100 pct
FROM
	(
	SELECT
		COUNT(*) CNT_BLOCK
	FROM
		V$SESSION
	WHERE
		BLOCKING_SESSION IS NOT NULL),
	(
	SELECT
		COUNT(*) CNT_ALL
	FROM
		V$SESSION)
UNION
SELECT
	'concurrency_rate',
	NVL(ROUND(SUM(duty_act.CNT * 100 / num_cores.VAL)), 0)
FROM
	(
		SELECT
			DECODE(SESSION_STATE, 'ON CPU', 'CPU', WAIT_CLASS) WAIT_CLASS, ROUND(COUNT(*) / (60 * 15), 1) CNT
		FROM
			V$ACTIVE_SESSION_HISTORY sh
		WHERE
			sh.SAMPLE_TIME >= SYSDATE - 15 / 1440
			AND DECODE(SESSION_STATE, 'ON CPU', 'CPU', WAIT_CLASS) IN ('Concurrency')
		GROUP BY
			DECODE(SESSION_STATE, 'ON CPU', 'CPU', WAIT_CLASS)
	) duty_act,
	(
		SELECT
			SUM(VALUE) VAL
		FROM
			V$OSSTAT
		WHERE
			STAT_NAME = 'NUM_CPU_CORES'
	) num_cores
	`

	rows, err := conn.Query(ctx, _sql, params["LockMaxTime"])
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

	// return format
	// {"active_background":49,"active_user":1,"concurrency_rate":0,"inactive_user":24,"lock_rate":0,"long_time_locked":0,"total":74}

	return strings.TrimSpace(string(jsonRes)), nil
}
