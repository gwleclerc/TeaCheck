package server

import (
	"crypto/tls"
	"strconv"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gwleclerc/TeaCheck/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func Serve(config config.Config) {
	engine := echo.New()
	engine.HideBanner = true
	engine.HidePort = true

	engine.Use(recoverMiddleware(), loggerMiddleware(), middleware.Gzip())

	log.WithField("port", config.ConfigListenPort).Info("Starting server")
	engine.Server.Addr = ":" + strconv.Itoa(config.ConfigListenPort)

	if config.TLSEnable {
		certificate, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
		if err != nil {
			log.WithFields(log.Fields{
				"tls-cert-file": config.TLSCertFile,
				"tls-key-file":  config.TLSKeyFile,
			}).Fatalf("Invalid certificate: %v", err)
		}
		engine.Server.TLSConfig = &tls.Config{
			NextProtos:   []string{"http/1.1"},
			Certificates: []tls.Certificate{certificate},
		}
	}

	if err := gracehttp.Serve(engine.Server); err != nil {
		log.Fatal(err)
	}
	log.Info("Shutting down gracefully")
}
