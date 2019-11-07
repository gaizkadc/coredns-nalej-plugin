/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
