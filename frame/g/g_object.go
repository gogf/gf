// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package g

import (
	"github.com/jin502437344/gf/database/gdb"
	"github.com/jin502437344/gf/database/gredis"
	"github.com/jin502437344/gf/frame/gins"
	"github.com/jin502437344/gf/i18n/gi18n"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/net/gtcp"
	"github.com/jin502437344/gf/net/gudp"
	"github.com/jin502437344/gf/os/gcfg"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gres"
	"github.com/jin502437344/gf/os/gview"
)

// Client is a convenience function, that creates and returns a new HTTP client.
func Client() *ghttp.Client {
	return ghttp.NewClient()
}

// Server returns an instance of http server with specified name.
func Server(name ...interface{}) *ghttp.Server {
	return gins.Server(name...)
}

// TCPServer returns an instance of tcp server with specified name.
func TCPServer(name ...interface{}) *gtcp.Server {
	return gtcp.GetServer(name...)
}

// UDPServer returns an instance of udp server with specified name.
func UDPServer(name ...interface{}) *gudp.Server {
	return gudp.GetServer(name...)
}

// View returns an instance of template engine object with specified name.
func View(name ...string) *gview.View {
	return gins.View(name...)
}

// Config returns an instance of config object with specified name.
func Config(name ...string) *gcfg.Config {
	return gins.Config(name...)
}

// Cfg is alias of Config.
// See Config.
func Cfg(name ...string) *gcfg.Config {
	return Config(name...)
}

// Resource returns an instance of Resource.
// The parameter <name> is the name for the instance.
func Resource(name ...string) *gres.Resource {
	return gins.Resource(name...)
}

// I18n returns an instance of gi18n.Manager.
// The parameter <name> is the name for the instance.
func I18n(name ...string) *gi18n.Manager {
	return gins.I18n(name...)
}

// Res is alias of Resource.
// See Resource.
func Res(name ...string) *gres.Resource {
	return Resource(name...)
}

// Log returns an instance of glog.Logger.
// The parameter <name> is the name for the instance.
func Log(name ...string) *glog.Logger {
	return gins.Log(name...)
}

// Database returns an instance of database ORM object with specified configuration group name.
func Database(name ...string) gdb.DB {
	return gins.Database(name...)
}

// DB is alias of Database.
// See Database.
func DB(name ...string) gdb.DB {
	return gins.Database(name...)
}

// Table is alias of Model.
func Table(tables string, db ...string) *gdb.Model {
	return DB(db...).Table(tables)
}

// Model creates and returns a model from specified database or default database configuration.
// The optional parameter <db> specifies the configuration group name of the database,
// which is "default" in default.
func Model(tables string, db ...string) *gdb.Model {
	return DB(db...).Model(tables)
}

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gredis.Redis {
	return gins.Redis(name...)
}
