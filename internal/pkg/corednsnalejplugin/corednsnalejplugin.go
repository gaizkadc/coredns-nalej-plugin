// Package log implements basic but useful request (access) logging plugin.
package corednsnalejplugin

import (
	"context"
	"errors"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/nalej/grpc-application-go"
	"github.com/rs/zerolog/log"
	"github.com/nalej/grpc-utils/pkg/conversions"
)

var (
	ErrNoAvailableEndpoints = errors.New("nalejPluginClient: no available endpoints")
	ErrOldCluster           = errors.New("nalejPluginClient: old cluster version")
)

type NalejPlugin struct {
	Next       plugin.Handler
	Fall       fall.F
	Zones      []string
	Upstream   *upstream.Upstream

	// SystemModelAddress with the host:port to connect to System Model
	SystemModelAddress string

	SMClient   grpc_application_go.ApplicationsClient
	Ctx        context.Context

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
			log.Info().Interface("question", question.Name).Msg("Incomming request")
			newRecords, err := np.ResolveEndpoint(question.Name)
			if err != nil {
				log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Str("question", question.Name).Msg("cannot resolve endpoint")
			}
			records = append(records, newRecords...)
		}

	}else{
		log.Error().Interface("state", state).Msg("unsupported query type")
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = append(m.Answer, records...)
	//m.Extra = append(m.Extra, extra...)

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func (np NalejPlugin) ResolveEndpoint(request string) ([]dns.RR, error) {
	// Query System Model
	smRequest := &grpc_application_go.GetAppEndPointRequest{
		Fqdn:                 request,
	}
	result, err := np.SMClient.GetAppEndpoints(context.Background(), smRequest)
	if err != nil{
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot retrieve endpoints from system model")
		return nil, err
	}
	log.Info().Int("len", len(result.AppEndpoints)).Msg("endpoints obtained")
	records := make([]dns.RR, 0)
	//for _, ep := range result.AppEndpoints{
	//	toAdd := &dns.A{Hdr: dns.RR_Header{Name: request, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 500}, A: net.IPv4(8, 8 ,8, 8)}
	//	log.Info().Interface("fqdn", ep.EndpointInstance.Fqdn).Msg("FQDN queried")
	//	records = append(records, toAdd)
	//}
	for _, ep := range result.AppEndpoints{
		toAdd := &dns.CNAME{Hdr: dns.RR_Header{Name: request, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 500}, Target: dns.Fqdn(ep.EndpointInstance.Fqdn)}
		records = append(records, toAdd)
	}

	return records, nil
}

// Name implements the Handler interface.
func (np NalejPlugin) Name() string { return "corednsnalejplugin" }

