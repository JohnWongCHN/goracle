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

	// return format
	// {"Buffer Cache Hit Ratio":100,"Memory Sorts Ratio":100,"Redo Allocation Hit Ratio":100,"User Transaction Per Sec":0.033,"Physical Reads Per Sec":0,"Physical Reads Per Txn":0,"Physical Writes Per Sec":0.317,"Physical Writes Per Txn":9.5,"Physical Reads Direct Per Sec":0,"Physical Reads Direct Per Txn":0,"Physical Writes Direct Per Sec":0,"Physical Writes Direct Per Txn":0,"Physical Reads Direct Lobs Per Sec":0,"Physical Reads Direct Lobs Per Txn":0,"Physical Writes Direct Lobs Per Sec":0,"Physical Writes Direct Lobs Per Txn":0,"Redo Generated Per Sec":205,"Redo Generated Per Txn":6150,"Logons Per Sec":0.017,"Logons Per Txn":0.5,"Open Cursors Per Sec":1.217,"Open Cursors Per Txn":36.5,"User Commits Per Sec":0,"User Commits Percentage":0,"User Rollbacks Per Sec":0.033,"User Rollbacks Percentage":100,"User Calls Per Sec":2.067,"User Calls Per Txn":62,"Recursive Calls Per Sec":14.933,"Recursive Calls Per Txn":448,"Logical Reads Per Sec":20.333,"Logical Reads Per Txn":610,"DBWR Checkpoints Per Sec":0,"Background Checkpoints Per Sec":0,"Redo Writes Per Sec":0.117,"Redo Writes Per Txn":3.5,"Long Table Scans Per Sec":0,"Long Table Scans Per Txn":0,"Total Table Scans Per Sec":1.1,"Total Table Scans Per Txn":33,"Full Index Scans Per Sec":0,"Full Index Scans Per Txn":0,"Total Index Scans Per Sec":0.367,"Total Index Scans Per Txn":11,"Total Parse Count Per Sec":1.267,"Total Parse Count Per Txn":38,"Hard Parse Count Per Sec":0.033,"Hard Parse Count Per Txn":1,"Parse Failure Count Per Sec":0.033,"Parse Failure Count Per Txn":1,"Cursor Cache Hit Ratio":29.73,"Disk Sort Per Sec":0,"Disk Sort Per Txn":0,"Rows Per Sort":21.559,"Execute Without Parse Ratio":38.71,"Soft Parse Ratio":97.368,"User Calls Ratio":12.157,"Host CPU Utilization (%)":2.132,"Network Traffic Volume Per Sec":4019.967,"Enqueue Timeouts Per Sec":0,"Enqueue Timeouts Per Txn":0,"Enqueue Waits Per Sec":0,"Enqueue Waits Per Txn":0,"Enqueue Deadlocks Per Sec":0,"Enqueue Deadlocks Per Txn":0,"Enqueue Requests Per Sec":217.65,"Enqueue Requests Per Txn":6529.5,"DB Block Gets Per Sec":0.867,"DB Block Gets Per Txn":26,"Consistent Read Gets Per Sec":19.467,"Consistent Read Gets Per Txn":584,"DB Block Changes Per Sec":0.783,"DB Block Changes Per Txn":23.5,"Consistent Read Changes Per Sec":0,"Consistent Read Changes Per Txn":0,"CPU Usage Per Sec":0.313,"CPU Usage Per Txn":9.394,"CR Blocks Created Per Sec":0,"CR Blocks Created Per Txn":0,"CR Undo Records Applied Per Sec":0,"CR Undo Records Applied Per Txn":0,"User Rollback UndoRec Applied Per Sec":0,"User Rollback Undo Records Applied Per Txn":0,"Leaf Node Splits Per Sec":0,"Leaf Node Splits Per Txn":0,"Branch Node Splits Per Sec":0,"Branch Node Splits Per Txn":0,"PX downgraded 1 to 25% Per Sec":0,"PX downgraded 25 to 50% Per Sec":0,"PX downgraded 50 to 75% Per Sec":0,"PX downgraded 75 to 99% Per Sec":0,"PX downgraded to serial Per Sec":0,"Physical Read Total IO Requests Per Sec":2.983,"Physical Read Total Bytes Per Sec":48878.933,"GC CR Block Received Per Second":0,"GC CR Block Received Per Txn":0,"GC Current Block Received Per Second":0,"GC Current Block Received Per Txn":0,"Global Cache Average CR Get Time":0,"Global Cache Average Current Get Time":0,"Physical Write Total IO Requests Per Sec":1.15,"Global Cache Blocks Corrupted":0,"Global Cache Blocks Lost":0,"Current Logons Count":81,"Current Open Cursors Count":2709,"User Limit %":0,"SQL Service Response Time":0.018,"Database Wait Time Ratio":0,"Database CPU Time Ratio":104.726,"Response Time Per Txn":8.97,"Row Cache Hit Ratio":100,"Row Cache Miss Ratio":0,"Library Cache Hit Ratio":99.482,"Library Cache Miss Ratio":0.518,"Shared Pool Free %":16.35,"PGA Cache Hit %":92.292,"Process Limit %":30,"Session Limit %":22.246,"Executions Per Txn":62,"Executions Per Sec":2.067,"Txns Per Logon":2,"Database Time Per Sec":0.299,"Physical Write Total Bytes Per Sec":15948.8,"Physical Read IO Requests Per Sec":0,"Physical Read Bytes Per Sec":0,"Physical Write IO Requests Per Sec":0.233,"Physical Write Bytes Per Sec":2594.133,"DB Block Changes Per User Call":0.379,"DB Block Gets Per User Call":0.419,"Executions Per User Call":1,"Logical Reads Per User Call":9.839,"Total Sorts Per User Call":0.274,"Total Table Scans Per User Call":0.532,"Current OS Load":0.17,"Streams Pool Usage Percentage":0,"PQ QC Session Count":0,"PQ Slave Session Count":0,"Queries parallelized Per Sec":0,"DML statements parallelized Per Sec":0,"DDL statements parallelized Per Sec":0,"PX operations not downgraded Per Sec":0,"Session Count":105,"Average Synchronous Single-Block Read Latency":0,"I/O Megabytes per Second":0.067,"I/O Requests per Second":4.133,"Average Active Sessions":0.003,"Active Serial Sessions":1,"Active Parallel Sessions":0,"Captured user calls":0,"Replayed user calls":0,"Workload Capture and Replay status":0,"Background CPU Usage Per Sec":1.303,"Background Time Per Sec":0.014,"Host CPU Usage Per Sec":4.183,"Cell Physical IO Interconnect Bytes":3889664,"Temp Space Used":5242880,"Total PGA Allocated":1265901568,"Total PGA Used by SQL Workareas":0,"Run Queue Per Sec":0,"VM in bytes Per Sec":0,"VM out bytes Per Sec":0}

	return strings.TrimSpace(string(jsonRes)), nil
}
