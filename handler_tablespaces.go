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

func tablespacesHandler(ctx context.Context, conn OraClient, params map[string]string,
	_ ...string) (interface{}, error) {

	_sql := `
SELECT  s.tablespace                                 as "tablespace_name",
		t.contents                                   as "contents",
		s.allocated_bytes                            as "file_bytes",
		s.max_bytes                                  as "max_bytes",
		s.free_bytes                                 as "free_bytes",
		s.used_bytes                                 as "used_bytes",
		s.used_pct_max                               as "used_pct_max",
		s.used_pct_allocated                         as "used_file_pct",
		s.free_pct_allocated                         as "free_file_pct",
		decode(t.status, 'ONLINE', 1, 'OFFLINE', 2,
				'READ ONLY', 3, 0)                   as "status"
FROM (
		SELECT  a.tablespace_name AS                                                tablespace,
				a.bytes_alloc                                                       allocated_bytes,
				nvl(b.bytes_free, 0)                                                free_bytes,
				a.bytes_alloc - nvl(b.bytes_free, 0)                                used_bytes,
				round((nvl(b.bytes_free, 0) / a.bytes_alloc) * 100, 2)              free_pct_allocated,
				100 - round((nvl(b.bytes_free, 0) / a.bytes_alloc) * 100, 2)        used_pct_allocated,
				maxbytes                                                            max_bytes,
				round(((a.bytes_alloc - nvl(b.bytes_free, 0)) / maxbytes) * 100, 2) used_pct_max
		FROM (
				SELECT f.tablespace_name,
						SUM(f.bytes)                                                    bytes_alloc,
						SUM(decode(f.autoextensible, 'YES', f.maxbytes, 'NO', f.bytes)) maxbytes
				FROM dba_data_files f
				GROUP BY tablespace_name
			) a,
			(
				SELECT f.tablespace_name,
						SUM(f.bytes) bytes_free
				FROM dba_free_space f
				GROUP BY tablespace_name
			) b
		WHERE a.tablespace_name = b.tablespace_name (+)
		UNION ALL
		SELECT  h.tablespace_name AS                                      tablespace,
				SUM(h.bytes_free + h.bytes_used)                          allocated_bytes,
				SUM((h.bytes_free + h.bytes_used) - nvl(p.bytes_used, 0)) free_bytes,
				SUM(nvl(p.bytes_used, 0))                                 used_bytes,
				round((SUM((h.bytes_free + h.bytes_used) - nvl(p.bytes_used, 0)) / SUM(h.bytes_used + h.bytes_free)) *
					100, 2)                                               free_pct_allocated,
				100 - round((SUM((h.bytes_free + h.bytes_used) - nvl(p.bytes_used, 0)) /
							SUM(h.bytes_used + h.bytes_free)) * 100, 2)  used_pct_allocated,
				SUM(f.maxbytes)                                           max_bytes,
				round(((SUM(nvl(p.bytes_used, 0))) /
					decode(SUM(f.maxbytes), 0, SUM(h.bytes_free + h.bytes_used), SUM(f.maxbytes))) * 100,
					2)                                                    used_pct_max
		FROM (
				SELECT DISTINCT *
				FROM sys.gv_$temp_space_header
			) h,
			(
				SELECT DISTINCT *
				FROM sys.gv_$temp_extent_pool
			) p,
			dba_temp_files f
		WHERE p.file_id (+) = h.file_id
		AND p.tablespace_name (+) = h.tablespace_name
		AND f.file_id = h.file_id
		AND f.tablespace_name = h.tablespace_name
		GROUP BY h.tablespace_name
	) s
		LEFT JOIN dba_tablespaces t ON s.tablespace = t.tablespace_name
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

	// return format
	// [{"USERS":{"contents":"PERMANENT","file_bytes":85854519296,"max_bytes":103081295872,"free_bytes":52018020352,"used_bytes":33836498944,"used_pct_max":83.29,"used_file_pct":39.41,"status":1}},{"UNDOTBS1":{"contents":"UNDO","file_bytes":959447040,"max_bytes":34359721984,"free_bytes":937099264,"used_bytes":22347776,"used_pct_max":2.79,"used_file_pct":2.33,"status":1}},{"SYSTEM":{"contents":"PERMANENT","file_bytes":2147483648,"max_bytes":34359721984,"free_bytes":991494144,"used_bytes":1155989504,"used_pct_max":6.25,"used_file_pct":53.83,"status":1}},{"SYSAUX":{"contents":"PERMANENT","file_bytes":4041211904,"max_bytes":34359721984,"free_bytes":231735296,"used_bytes":3809476608,"used_pct_max":11.76,"used_file_pct":94.27,"status":1}},{"INDX":{"contents":"PERMANENT","file_bytes":30482104320,"max_bytes":51547979776,"free_bytes":4262002688,"used_bytes":26220101632,"used_pct_max":59.13,"used_file_pct":86.02,"status":1}},{"JZX_COMMT":{"contents":"PERMANENT","file_bytes":4297064448,"max_bytes":4297064448,"free_bytes":4296015872,"used_bytes":1048576,"used_pct_max":100,"used_file_pct":0.02,"status":1}},{"TEMP":{"contents":"TEMPORARY","file_bytes":3276800000,"max_bytes":34359721984,"free_bytes":3271557120,"used_bytes":5242880,"used_pct_max":9.54,"used_file_pct":0.16,"status":1}}]

	return "[" + strings.Join(data, ",") + "]", nil
}
