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

	"git.zabbix.com/ap/plugin-support/metric"
	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/uri"
)

const (
	keyArchive                = "goracle.archive.info"
	keyArchiveDiscovery       = "goracle.archive.discovery"
	keyASMDiskGroups          = "goracle.diskgroups.stats"
	keyASMDiskGroupsDiscovery = "goracle.diskgroups.discovery"
	keyCDB                    = "goracle.cdb.info"
	keyCustomQuery            = "goracle.custom.query"
	keyDatabasesDiscovery     = "goracle.db.discovery"
	keyDataFiles              = "goracle.datafiles.stats"
	keyFRA                    = "goracle.fra.stats"
	keyInstance               = "goracle.instance.info"
	keyPDB                    = "goracle.pdb.info"
	keyPDBDiscovery           = "goracle.pdb.discovery"
	keyPGA                    = "goracle.pga.stats"
	keyPing                   = "goracle.ping"
	keyProc                   = "goracle.proc.stats"
	keyRedoLog                = "goracle.redolog.info"
	keySessions               = "goracle.sessions.stats"
	keySGA                    = "goracle.sga.stats"
	keySysMetrics             = "goracle.sys.metrics"
	keySysParams              = "goracle.sys.params"
	keyTablespaces            = "goracle.ts.stats"
	keyTablespacesDiscovery   = "goracle.ts.discovery"
	keyUser                   = "goracle.user.info"
	keyDatabaseSegment        = "goracle.db.segment_size"
)

// handlerFunc defines an interface must be implemented by handlers.
type handlerFunc func(ctx context.Context, conn OraClient,
	params map[string]string, extraParams ...string) (res interface{}, err error)

// getHandlerFunc returns a handlerFunc related to a given key.
func getHandlerFunc(key string) handlerFunc {
	switch key {
	case keyASMDiskGroups:
		return asmDiskGroupsHandler
	case keyASMDiskGroupsDiscovery:
		return asmDiskGroupsDiscovery
	case keyArchive:
		return archiveHandler
	case keyArchiveDiscovery:
		return archiveDiscoveryHandler
	case keyCDB:
		return cdbHandler
	case keyCustomQuery:
		return customQueryHandler
	case keyDataFiles:
		return dataFileHandler
	case keyDatabasesDiscovery:
		return databasesDiscoveryHandler
	case keyFRA:
		return fraHandler
	case keyInstance:
		return instanceHandler
	case keyPDB:
		return pdbHandler
	case keyPDBDiscovery:
		return pdbDiscoveryHandler
	case keyPGA:
		return pgaHandler
	case keyPing:
		return pingHandler
	case keyProc:
		return procHandler
	case keyRedoLog:
		return redoLogHandler
	case keySGA:
		return sgaHandler
	case keySessions:
		return sessionsHandler
	case keySysMetrics:
		return sysMetricsHandler
	case keySysParams:
		return sysParamsHandler
	case keyTablespaces:
		return tablespacesHandler
	case keyTablespacesDiscovery:
		return tablespacesDiscoveryHandler
	case keyUser:
		return userHandler
	case keyDatabaseSegment:
		return databaseSegmentHandler

	default:
		return nil
	}
}

var uriDefaults = &uri.Defaults{Scheme: "tcp", Port: "1521"}

// Common params: [URI|Session][,User][,Password][,Service]
var (
	paramURI = metric.NewConnParam("URI", "URI to connect or session name.").
			WithDefault(uriDefaults.Scheme + "://localhost:" + uriDefaults.Port).WithSession().
			WithValidator(uri.URIValidator{Defaults: uriDefaults, AllowedSchemes: []string{"tcp"}})
	paramUsername = metric.NewConnParam("User", "Oracle user.").WithDefault("")
	paramPassword = metric.NewConnParam("Password", "User's password.").WithDefault("")
	paramService  = metric.NewConnParam("Service", "Service name to be used for connection.").
			WithDefault("XE")
)

var metrics = metric.MetricSet{
	keyASMDiskGroups: metric.New("Returns ASM disk groups statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyASMDiskGroupsDiscovery: metric.New("Returns list of ASM disk groups in LLD format.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyArchive: metric.New("Returns archive logs statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyArchiveDiscovery: metric.New("Returns list of archive logs in LLD format.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyCDB: metric.New("Returns CDBs info.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyCustomQuery: metric.New("Returns result of a custom query.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService,
			metric.NewParam("QueryName", "Name of a custom query "+
				"(must be equal to a name of an SQL file without an extension).").SetRequired(),
		}, true),

	keyDataFiles: metric.New("Returns data files statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyDatabasesDiscovery: metric.New("Returns list of databases in LLD format.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyFRA: metric.New("Returns FRA statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyInstance: metric.New("Returns instance stats.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyPDB: metric.New("Returns PDBs info.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyPDBDiscovery: metric.New("Returns list of PDBs in LLD format.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyPGA: metric.New("Returns PGA statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyPing: metric.New("Tests if connection is alive or not.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyProc: metric.New("Returns processes statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyRedoLog: metric.New("Returns log file information from the control file.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keySGA: metric.New("Returns SGA statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keySessions: metric.New("Returns sessions statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService,
			metric.NewParam("LockMaxTime", "Maximum session lock duration in seconds to count "+
				"the session as a prolongedly locked.").WithDefault("600").WithValidator(metric.NumberValidator{}),
		}, false),

	keySysMetrics: metric.New("Returns a set of system metric values.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService,
			metric.NewParam("Duration", "Capturing interval in seconds of system metric values.").
				WithDefault("60").WithValidator(metric.SetValidator{Set: []string{"15", "60"}}),
		}, false),

	keySysParams: metric.New("Returns a set of system parameter values.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyTablespaces: metric.New("Returns tablespaces statistics.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyTablespacesDiscovery: metric.New("Returns list of tablespaces in LLD format.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),

	keyUser: metric.New("Returns user information.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService,
			metric.NewParam("Username", "Username for which the information is needed."),
		}, false),
	keyDatabaseSegment: metric.New("Returns database segment size.",
		[]*metric.Param{paramURI, paramUsername, paramPassword, paramService}, false),
}

func init() {
	plugin.RegisterMetrics(&impl, pluginName, metrics.List()...)
}
