package updater

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	config2 "free/config"
	"strings"
)

func parseV2ray(raw string) []config2.V2rayConfig {
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil
	}
	var configs []config2.V2rayConfig
	urls := strings.Split(string(data), "\n")
	for _, u := range urls {
		c := parseSingleV2ray(u)
		if c.Validate() {
			configs = append(configs)
		}
	}
	return configs
}

func parseSingleV2ray(raw string) config2.V2rayConfig {
	body := strings.TrimPrefix(raw, "vmess://")
	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return config2.V2rayConfig{}
	}
	var mp = make(map[string]string)
	err = json.Unmarshal(decodedBody, &mp)
	if err != nil {
		return config2.V2rayConfig{}
	}

	// 重写
	mp["ps"] = fmt.Sprintf("%s_%s_%s", fqFree, mp["host"], mp["port"])
	data, err := json.Marshal(mp)
	if err != nil {
		return config2.V2rayConfig{}
	}
	config := config2.V2rayConfig{
		Host: mp["host"],
		Net:  mp["net"],
		Port: mp["port"],
		Raw:  fmt.Sprintf("vmess://%s", string(data)),
	}
	return config
}

const (
	fqFree = "fqfree.xyz"
)
