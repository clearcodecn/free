package updater

import (
	"free/config"
	"strings"
)

/**
// aes-256-gcm:8n6pwAcrrv2pj6tFY2p3TbQ6
YWVzLTI1Ni1nY206OG42cHdBY3JydjJwajZ0RlkycDNUYlE2
@104.200.131.165:33992

ss://YWVzLTI1Ni1nY206OG42cHdBY3JydjJwajZ0RlkycDNUYlE2@104.200.131.165:33992#github.com/freefq%20-%20%E5%8C%97%E7%BE%8E%E5%9C%B0%E5%8C%BA%20%2018

*/
func parseSingleSS(s string) config.SSConfig {
	withoutSuffix := strings.SplitN(s, "#", 2)[0]
	x := strings.TrimPrefix(withoutSuffix, "ss://")

	parts := strings.SplitN(x, "@", 2)
	hostPort := strings.Split(parts[1], ":")
	if len(hostPort) == 2 {
		host, port := hostPort[0], hostPort[1]
		return config.SSConfig{
			Host: host,
			Port: port,
			Raw:  withoutSuffix,
		}
	}
	return config.SSConfig{}
}
