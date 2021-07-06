package updater

import (
	"encoding/base64"
	"free/config"
	"net/url"
	"strings"
)

func parseSSR(data []byte) []config.SSRConfig {
	s, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil
	}
	var configs []config.SSRConfig
	arr := strings.Split(string(s), "\n")
	for _, a := range arr {
		c, ok := parseSingleSSR(a)
		if ok {
			configs = append(configs, c)
		}
	}
	return configs
}

func parseSingleSSR(s string) (config.SSRConfig, bool) {
	old := s
	s = strings.TrimPrefix(s, "ssr://")
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)
	s = fill(s)
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return config.SSRConfig{}, false
	}
	c, ok := raw2Config(string(raw))
	if !ok {
		return c, false
	}
	c.Raw = old
	return c, true
}

func fill(s string) string {
	l := 4 - len(s)%4
	if l > 0 && l != 4 {
		s = s + strings.Repeat("=", l)
	}
	return s
}

// host:port:protocol:method:obfs:base64pass/?obfsparam=base64param&protoparam=base64param&remarks=base64remarks&group=base64group&udpport=0&uot=0
// ssr://server:port:protocol:method:obfs:password_base64/?params_base64
func raw2Config(raw string) (config.SSRConfig, bool) {
	parts := strings.Split(raw, "/?")

	var (
		paths string
		query string
	)
	paths = parts[0]
	if len(parts) > 1 {
		query = parts[1]
	}

	arr := strings.Split(paths, ":")
	if len(arr) != 6 {
		return config.SSRConfig{}, false
	}
	ip, port, protocol, method, obfs, passwordb64 := arr[0], arr[1], arr[2], arr[3], arr[4], arr[5]
	// parse query
	var params = make(map[string]string)
	if len(query) != 0 {
		q, err := url.ParseQuery(query)
		if err != nil {
			return config.SSRConfig{}, false
		}
		for k := range q {
			params[k] = q.Get(k)
		}
	}

	c := config.SSRConfig{
		Server:   ip,
		Port:     port,
		Protocol: protocol,
		Method:   method,
		Obfs:     obfs,
		Password: passwordb64,
		Params:   params,
	}

	if !c.Validate() {
		return c, false
	}

	return c, true
}
