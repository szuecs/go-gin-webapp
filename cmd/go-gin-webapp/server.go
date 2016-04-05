package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/golang/glog"
	"github.com/szuecs/go-gin-webapp/conf"
	"github.com/szuecs/go-gin-webapp/frontend"
	"golang.org/x/oauth2"
)

// Buildstamp is used to store the timestamp of the build
var Buildstamp = "Not set"

// Githash is used to store the commit hash of the build
var Githash = "Not set"

// Version is used to store the tagged version of the build
var Version = "Not set"

// flag variables
var version bool

var serverConfig *conf.Config

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
Example:
  %% %s
`, bin, bin)
		flag.PrintDefaults()
	}

	var err error
	serverConfig, err = conf.New()
	if err != nil {
		fmt.Printf("Could not read config, caused by: %s\n", err)
		os.Exit(2)
	}

	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&serverConfig.DebugEnabled, "debug", serverConfig.DebugEnabled, "Enable debug output")
	flag.BoolVar(&serverConfig.Oauth2Enabled, "oauth", serverConfig.Oauth2Enabled, "Enable OAuth2")
	flag.BoolVar(&serverConfig.ProfilingEnabled, "profile", serverConfig.ProfilingEnabled, "Enable profiling.")
	flag.StringVar(&serverConfig.AuthURL, "oauth-authurl", serverConfig.AuthURL, "OAuth2 Auth URL")
	flag.StringVar(&serverConfig.TokenURL, "oauth-tokeninfourl", serverConfig.TokenURL, "OAuth2 Auth URL")
	flag.StringVar(&serverConfig.TLSCertfilePath, "tls-cert", serverConfig.TLSCertfilePath, "TLS Certfile")
	flag.StringVar(&serverConfig.TLSKeyfilePath, "tls-key", serverConfig.TLSKeyfilePath, "TLS Keyfile")
	flag.IntVar(&serverConfig.Port, "port", serverConfig.Port, "Listening TCP Port of the service.")
	flag.IntVar(&serverConfig.MonitorPort, "monitor-port", serverConfig.MonitorPort, "Listening TCP Port of the monitor.")
	flag.DurationVar(&serverConfig.LogFlushInterval, "flush-interval", serverConfig.LogFlushInterval, "Interval to flush Logs to disk.")

	flag.Parse()
}

func main() {
	if version {
		fmt.Printf(`%s Version: %s
================================
    Buildtime: %s
    GitHash: %s
`, path.Base(os.Args[0]), Version, Buildstamp, Githash)
		os.Exit(0)
	}

	httpOnly := IsHTTPOnly(serverConfig)
	var cfg *frontend.ServiceConfig
	cfg, err := GetServiceConfig(serverConfig, httpOnly)

	svc := frontend.Service{
		Healthy: false,
	}
	err = svc.Run(cfg)
	if err != nil {
		fmt.Printf("ERR: Could not start service, caused by: %s\n", err)
		os.Exit(1)
	}
}

// IsHTTPOnly defaults to false and returns true if it can not read Cert or Key.
func IsHTTPOnly(cfg *conf.Config) bool {
	var err error
	if _, err = os.Stat(cfg.TLSCertfilePath); os.IsNotExist(err) {
		glog.Warningf("WARN: No Certfile found %s\n", cfg.TLSCertfilePath)
		return true
	} else if _, err = os.Stat(cfg.TLSKeyfilePath); os.IsNotExist(err) {
		glog.Warningf("WARN: No Keyfile found %s\n", cfg.TLSKeyfilePath)
		return true
	}
	return false
}

// GetServiceConfig returns frontend.ServiceConfig. err is nil unless
// tls.LoadX509KeyPair failed to load configured Cert or Key.
func GetServiceConfig(cfg *conf.Config, httpOnly bool) (*frontend.ServiceConfig, error) {
	var keypair tls.Certificate
	var err error
	if !httpOnly {
		keypair, err = tls.LoadX509KeyPair(cfg.TLSCertfilePath, cfg.TLSKeyfilePath)
		if err != nil {
			return nil, fmt.Errorf("ERR: Could not load X509 KeyPair, caused by: %s\n", err)

		}
	}

	var oauth2Endpoint = oauth2.Endpoint{
		AuthURL:  cfg.AuthURL,
		TokenURL: cfg.TokenURL,
	}

	return &frontend.ServiceConfig{
		Config:          cfg,
		OAuth2Endpoints: oauth2Endpoint,
		CertKeyPair:     keypair,
		Httponly:        httpOnly,
	}, nil
}
