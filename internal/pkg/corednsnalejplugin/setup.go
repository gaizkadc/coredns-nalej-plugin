package corednsnalejplugin

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
	"github.com/nalej/grpc-application-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"strings"
)

func init() {
	fmt.Println("plugin.init")
	caddy.RegisterPlugin("corednsnalejplugin", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	n, err := corednsnalejpluginParse(c)
	if err != nil {
		return plugin.Error("corednsnalejplugin", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		n.Next = next
		return n
	})

	return nil
}

func corednsnalejpluginParse (c *caddy.Controller) (*NalejPlugin, error) {
	nalejPlugin := NalejPlugin{
		Ctx:        context.Background(),
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	for c.Next() {
		nalejPlugin.Zones = c.RemainingArgs()
		if len(nalejPlugin.Zones) == 0 {
			nalejPlugin.Zones = make([]string, len(c.ServerBlockKeys))
			copy(nalejPlugin.Zones, c.ServerBlockKeys)
		}
		for i, str := range nalejPlugin.Zones {
			nalejPlugin.Zones[i] = plugin.Host(str).Normalize()
		}

		if c.NextBlock() {
			for {
				switch c.Val() {

				case "systemModelAddress":
					address := c.RemainingArgs()
					if len(address) != 1 {
						return &NalejPlugin{}, c.Errf("system model address expected")
					}
					nalejPlugin.SystemModelAddress = address[0]
				case "debug":
					zerolog.SetGlobalLevel(zerolog.DebugLevel)
				default:
					if c.Val() != "}" {
						return &NalejPlugin{}, c.Errf("unknown property '%s'", c.Val())
					}
				}

				if !c.Next() {
					break
				}
			}

		}
		log.Info().Str("URL", nalejPlugin.SystemModelAddress).Msg("System Model")

		sp := strings.Split(nalejPlugin.SystemModelAddress, ":")
		ips, err := net.LookupIP(sp[0])
		if err != nil {
			log.Error().Err(err).Msg("cannot get ips")
		}
		for _, ip := range ips {
			log.Info().Str("A", ip.String()).Msg("answer")
		}

		smConn, err := grpc.Dial(nalejPlugin.SystemModelAddress, grpc.WithInsecure())
		if err != nil{
			return nil, c.Errf("cannot create connection with system model")
		}

		nalejPlugin.SMClient = grpc_application_go.NewApplicationsClient(smConn)

		return &nalejPlugin, nil
	}
	return &NalejPlugin{}, nil
}
