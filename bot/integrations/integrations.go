package integrations

import (
	"github.com/TicketsBot-cloud/common/integrations/bloxlink"
	"github.com/TicketsBot-cloud/common/integrations/rover"
	"github.com/TicketsBot-cloud/common/webproxy"
	"github.com/TicketsBot-cloud/worker/bot/redis"
	"github.com/TicketsBot-cloud/worker/config"
)

var (
	WebProxy    *webproxy.WebProxy
	SecureProxy *SecureProxyClient
	Bloxlink    *bloxlink.BloxlinkIntegration
	Rover       *rover.RoverIntegration
)

func InitIntegrations() {
	WebProxy = webproxy.NewWebProxy(config.Conf.WebProxy.Url, config.Conf.WebProxy.AuthHeaderName, config.Conf.WebProxy.AuthHeaderValue)
	Bloxlink = bloxlink.NewBloxlinkIntegration(redis.Client, WebProxy, config.Conf.Integrations.BloxlinkApiKey)
	Rover = rover.NewRoverIntegration(redis.Client, WebProxy, config.Conf.Integrations.RoverApiKey)
	SecureProxy = NewSecureProxy(config.Conf.Integrations.SecureProxyUrl)
}
