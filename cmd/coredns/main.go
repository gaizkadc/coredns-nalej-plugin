/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/nalej/coredns-nalej-plugin/internal/pkg/corednsnalejplugin"
	"github.com/nalej/coredns-nalej-plugin/version"
)

var DebugLevel bool
var MainVersion string
var MainCommit string
var directives = []string{
	"corednsnalejplugin",
}

func init() {
	dnsserver.Directives = directives
}

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	coremain.Run()
}
