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

// Package log implements basic but useful request (access) logging plugin.
package corednsnalejplugin

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"net"
)

const defaultTTL = 5

type NalejPlugin struct {
	Next     plugin.Handler
	Fall     fall.F
	Zones    []string
	Upstream *upstream.Upstream

	// SystemModelAddress with the host:port to connect to System Model
	SystemModelAddress string

	SMClient grpc_application_go.ApplicationsClient
	Ctx      context.Context

	endpoints []string // Stored here as well, to aid in testing.
}

// ServeDNS implements the plugin.Handler interface.
func (np NalejPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	zone := plugin.Zones(np.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(np.Name(), np.Next, ctx, w, r)
	}

	var records []dns.RR

	if state.QType() == dns.TypeA {
		for _, question := range state.Req.Question {
			log.Debug().Interface("question", question.Name).Msg("Incomming request")
			newRecords, err := np.ResolveEndpoint(question.Name)
			if err != nil {
				log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Str("question", question.Name).Msg("cannot resolve endpoint")
			} else {
				records = append(records, newRecords...)
			}
		}

	} else {
		log.Error().Interface("state", state).Msg("unsupported query type")
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = append(m.Answer, records...)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

func (np NalejPlugin) ResolveEndpoint(request string) ([]dns.RR, error) {
	// Query System Model
	smRequest := &grpc_application_go.GetAppEndPointRequest{
		Fqdn: request,
	}
	result, err := np.SMClient.GetAppEndpoints(context.Background(), smRequest)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot retrieve endpoints from system model")
		return nil, err
	}
	log.Debug().Int("len", len(result.AppEndpoints)).Msg("endpoints obtained")
	records := make([]dns.RR, 0)

	for _, ep := range result.AppEndpoints {
		if ep.EndpointInstance.Type == grpc_application_go.EndpointType_WEB {
			toAdd := &dns.CNAME{
				Hdr:    dns.RR_Header{Name: request, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: defaultTTL},
				Target: dns.Fqdn(ep.EndpointInstance.Fqdn),
			}
			records = append(records, toAdd)
			log.Debug().Str("target", toAdd.Target).Msg("CNAME")
		} else if ep.EndpointInstance.Type == grpc_application_go.EndpointType_INGESTION {
			toAdd := &dns.A{
				Hdr: dns.RR_Header{Name: request, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: defaultTTL},
				A:   net.ParseIP(ep.EndpointInstance.Fqdn),
			}
			records = append(records, toAdd)
			log.Debug().Str("IP", toAdd.A.String()).Msg("A")
		} else {
			log.Warn().Interface("request type", ep.EndpointInstance.Type).Msg("request type not recognized")
		}
	}

	return records, nil
}

// Name implements the Handler interface.
func (np NalejPlugin) Name() string { return "corednsnalejplugin" }
