# Go based Zabbix Agent 2 external plugin for Oracle database

Provides native Zabbix solution for monitoring Oracle Database (multi-model database management system).
It can monitor several Oracle instances simultaneously, remote or local to the Zabbix Agent.
The plugin keeps connections in the open state to reduce network congestion, latency, CPU and
memory usage. Best for use in conjunction with the official
[Oracle template.](https://git.zabbix.com/projects/ZBX/repos/zabbix/browse/templates/db/oracle_agent2)
You can extend it or create your template for your specific needs.

## Requirements

* Zabbix Agent 2
* Go >= 1.18 (required only to build from source)
* Zig >= 0.9.1 (required only to build from source)

## Supported versions

* Oracle Database 10g
* Oracle Database 11g
* Oracle Database 12c2
* Oracle Database 18c
* Oracle Database 19c

## Installation

* Create an Oracle DB user and grant permissions

```sql
CREATE USER zabbix_mon IDENTIFIED BY <PASSWORD>;
-- Grant access to the zabbix_mon user.
GRANT CONNECT, CREATE SESSION TO zabbix_mon;
GRANT SELECT ON DBA_TABLESPACE_USAGE_METRICS TO zabbix_mon;
GRANT SELECT ON DBA_TABLESPACES TO zabbix_mon;
GRANT SELECT ON DBA_USERS TO zabbix_mon;
GRANT SELECT ON SYS.DBA_DATA_FILES TO zabbix_mon;
GRANT SELECT ON V$ACTIVE_SESSION_HISTORY TO zabbix_mon;
GRANT SELECT ON V$ARCHIVE_DEST TO zabbix_mon;
GRANT SELECT ON V$ASM_DISKGROUP TO zabbix_mon;
GRANT SELECT ON V$DATABASE TO zabbix_mon;
GRANT SELECT ON V$DATAFILE TO zabbix_mon;
GRANT SELECT ON V$INSTANCE TO zabbix_mon;
GRANT SELECT ON V$LOG TO zabbix_mon;
GRANT SELECT ON V$OSSTAT TO zabbix_mon;
GRANT SELECT ON V$PGASTAT TO zabbix_mon;
GRANT SELECT ON V$PROCESS TO zabbix_mon;
GRANT SELECT ON V$RECOVERY_FILE_DEST TO zabbix_mon;
GRANT SELECT ON V$RESTORE_POINT TO zabbix_mon;
GRANT SELECT ON V$SESSION TO zabbix_mon;
GRANT SELECT ON V$SGASTAT TO zabbix_mon;
GRANT SELECT ON V$SYSMETRIC TO zabbix_mon;
GRANT SELECT ON V$SYSTEM_PARAMETER TO zabbix_mon;
```

* Make sure a TNS Listener and an Oracle instance are available for connection.  

## Configuration

The Zabbix agent 2 configuration file is used to configure plugins.

**Plugins.Goracle.CallTimeout** — The maximum time in seconds for waiting when a request has to be done.  
*Default value:* equals the global Timeout configuration parameter.  
*Limits:* 1-30

**Plugins.Goracle.ConnectTimeout** — The maximum time in seconds for waiting when a connection has to be established.  
*Default value:* equals the global Timeout configuration parameter.  
*Limits:* 1-30

**Plugins.Goracle.CustomQueriesPath** — Full pathname of a directory containing *.sql* files with custom queries.  
*Default value:* — (the feature is disabled by default)

**Plugins.Goracle.KeepAlive** — Sets a time for waiting before unused connections will be closed.  
*Default value:* 300 sec.  
*Limits:* 60-900
