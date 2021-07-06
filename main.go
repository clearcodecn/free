package main

import (
	"free/pkg/updater"
	"os"
)

// ssr://server:port:protocol:method:obfs:password_base64/?params_base64
func main() {
	os.Setenv("https_proxy", "http://localhost:1087")
	updater.Update()
}
