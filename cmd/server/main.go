package main

import (
	"errors"
	"fmt"
	"free/config"
	"free/internal/http"
	"free/pkg/ip"
	"free/pkg/service"
	"free/pkg/updater"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	g        errgroup.Group
	stopChan = make(chan struct{})
)

func init() {
	config.InitConfig()
	dur, _ := time.ParseDuration(config.GetConfig().Updater.Duration)
	if dur == 0 {
		dur = time.Hour * 4
	}
	fmt.Println(fmt.Sprintf("%+v", config.GetConfig()))
	service.InitService(config.GetConfig().Updater.CachePath, dur)

	ip.InitIPDB(config.GetConfig().Updater.Ip2Region)
}

func main() {
	/**
	本地调试，加代理.
	*/
	os.Setenv("https_proxy", "socks5://localhost:1080")

	g.Go(func() error {
		if err := http.Serve(); err != nil {
			if errors.Is(err, http2.ErrServerClosed) {
				return nil
			}
			return err
		}
		return nil
	})
	g.Go(func() error {
		dur, _ := time.ParseDuration(config.GetConfig().Updater.Duration)
		if dur == 0 {
			dur = time.Hour * 4
		}

		updater.Do(stopChan, &updater.Config{
			Duration:       dur,
			LastRunTimFile: config.GetConfig().Updater.LastRunTimeFile,
			StorageFile:    config.GetConfig().Updater.CachePath,
		})

		return nil
	})
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			si := <-c
			switch si {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				return shutdown()
			case syscall.SIGHUP:
			default:
				return nil
			}
		}
	})

	if err := g.Wait(); err != nil {
		logrus.Error("服务运行失败: ", err)
		panic(err)
	}
}

func shutdown() error {
	close(stopChan)
	if err := http.Stop(); err != nil {
		return err
	}
	return nil
}
