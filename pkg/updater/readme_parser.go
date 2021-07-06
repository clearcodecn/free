package updater

import (
	"bufio"
	"bytes"
	"free/config"
	"strings"
)

func parseReadme(data []byte) ([]config.SSConfig, []config.SSRConfig, []config.V2rayConfig) {
	buf := bufio.NewReader(bytes.NewReader(data))
	var (
		ss    []string
		vmess []string
		ssr   []string
	)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		if strings.HasPrefix(string(line), "vmess://") {
			vmess = append(vmess, string(line))
		}
		if strings.HasPrefix(string(line), "ss://") {
			ss = append(ss, string(line))
		}
		if strings.HasPrefix(string(line), "ssr://") {
			ssr = append(ssr, string(line))
		}
	}

	var (
		ssConfig  []config.SSConfig
		ssrConfig []config.SSRConfig
		v2Config  []config.V2rayConfig
	)

	for _, v := range ss {
		x := parseSingleSS(v)
		if x.Host == "" {
			continue
		}
		ssConfig = append(ssConfig, x)
	}

	for _, v := range vmess {
		v2Config = append(v2Config, parseSingleV2ray(v))
	}
	for _, v := range ssr {
		x, ok := parseSingleSSR(v)
		if ok {
			ssrConfig = append(ssrConfig, x)
		}
	}

	return ssConfig, ssrConfig, v2Config
}
