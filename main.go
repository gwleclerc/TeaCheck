package main

import (
	"os"

	"github.com/gwleclerc/TeaCheck/config"
	"github.com/gwleclerc/TeaCheck/server"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

var appName, buildVersion, buildCommit, buildDate string // nolint

func parseConfig() (c config.Config) {
	c.Build = config.Build{
		AppName:      appName,
		BuildVersion: buildVersion,
		BuildCommit:  buildCommit,
		BuildDate:    buildDate,
	}

	// Use a prefix for environment variables
	flag.CommandLine = flag.NewFlagSetWithEnvPrefix(os.Args[0], "TEACHECK", flag.ExitOnError)
	flag.StringVar(&c.LogLevel, "log-level", "info", "Available levels: panic, fatal, error, warning, info, debug, trace")
	flag.StringVar(&c.ConfigBasePath, "config-base-path", "/", "Base path applied to TeaCheck UI")
	flag.IntVar(&c.ConfigListenPort, "config-listen-port", 8081, "Server listening port")
	flag.BoolVar(&c.TLSEnable, "tls-enable", false, "Enable TLS using the provided certificate")
	flag.StringVar(&c.TLSCertFile, "tls-cert-file", "/etc/smocker/tls/certs/cert.pem", "Path to TLS certificate file ")
	flag.StringVar(&c.TLSKeyFile, "tls-private-key-file", "/etc/smocker/tls/private/key.pem", "Path to TLS key file")
	flag.Parse()
	return
}

func setupLogger(logLevel string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithError(err).WithField("log-level", level).Warn("Invalid log level, fallback to info")
		level = log.InfoLevel
	}
	log.WithField("log-level", level).Info("Setting log level")
	log.SetLevel(level)
}

func main() {
	c := parseConfig()
	setupLogger(c.LogLevel)
	server.Serve(c)
}
