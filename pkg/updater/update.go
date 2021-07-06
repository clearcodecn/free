package updater

import (
	"free/config"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	v2URL = `https://raw.githubusercontent.com/freefq/free/master/v2`
	// ssr://server:port:protocol:method:obfs:password_base64/?params_base64
	ssrURL    = `https://raw.githubusercontent.com/freefq/free/master/ssr`
	readmeURL = `https://github.com/freefq/free/raw/master/README.md`
)

func Update() config.AllConfig {
	var (
		v2ray []config.V2rayConfig
		ssr   []config.SSRConfig
		ss    []config.SSConfig
	)
	v2data, ssrdata, readme := download()
	if len(ssrdata) != 0 {
		ssr = parseSSR(ssrdata)
	}
	if len(v2data) != 0 {
		v2ray = parseV2ray(string(v2data))
	}
	var (
		a []config.SSConfig
		b []config.SSRConfig
		c []config.V2rayConfig
	)
	if len(readme) != 0 {
		a, b, c = parseReadme(readme)
	}
	ssr = append(ssr, b...)
	v2ray = append(v2ray, c...)
	ss = append(ss, a...)

	return config.AllConfig{
		SS:    ss,
		SSR:   ssr,
		V2ray: v2ray,
	}
}

func download() ([]byte, []byte, []byte) {
	var (
		g          errgroup.Group
		v2Data     []byte
		ssrData    []byte
		readmeData []byte
	)
	g.Go(func() error {
		resp, err := http.Get(v2URL)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		v2Data = data
		return nil
	})
	g.Go(func() error {
		resp, err := http.Get(ssrURL)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		ssrData = data
		return nil
	})
	g.Go(func() error {
		resp, err := http.Get(readmeURL)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		readmeData = data
		return nil
	})
	if err := g.Wait(); err != nil {
		logrus.WithError(err).Errorf("failed to download")
	}
	return v2Data, ssrData, readmeData
}
