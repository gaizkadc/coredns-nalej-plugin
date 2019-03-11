/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package main

import (
    "fmt"
    "github.com/coredns/coredns/core/dnsserver"
    "github.com/coredns/coredns/coremain"
    "github.com/nalej/golang-template/version"
    _ "github.com/coredns/proxy"
    _ "github.com/coredns/coredns/plugin/whoami"
    //_ "github.com/coredns/coredns/plugin/log"
    _ "github.com/nalej/coredns-nalej-plugin/internal/pkg/corednsnalejplugin"
)

var MainVersion string
var MainCommit string
var directives = []string{
    "corednsnalejplugin",
    "proxy",
    "whoami",
    //"log",
}

func init() {
    dnsserver.Directives = directives
}

func main() {
    version.AppVersion = MainVersion
    version.Commit = MainCommit
    coremain.Run()
    fmt.Println("executed")
}