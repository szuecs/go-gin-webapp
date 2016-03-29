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

	fmt.Printf("serverConfig: %+v\n", serverConfig)

	// default https, if cert and key are found
	var err error
	httpOnly := false
	if _, err = os.Stat(serverConfig.TLSCertfilePath); os.IsNotExist(err) {
		glog.Warningf("WARN: No Certfile found %s\n", serverConfig.TLSCertfilePath)
		httpOnly = true
	} else if _, err = os.Stat(serverConfig.TLSKeyfilePath); os.IsNotExist(err) {
		glog.Warningf("WARN: No Keyfile found %s\n", serverConfig.TLSKeyfilePath)
		httpOnly = true
	}
	var keypair tls.Certificate
	if httpOnly {
		keypair = tls.Certificate{}
	} else {
		keypair, err = tls.LoadX509KeyPair(serverConfig.TLSCertfilePath, serverConfig.TLSKeyfilePath)
		if err != nil {
			fmt.Printf("ERR: Could not load X509 KeyPair, caused by: %s\n", err)
			os.Exit(1)
		}
	}

	var oauth2Endpoint = oauth2.Endpoint{
		AuthURL:  serverConfig.AuthURL,
		TokenURL: serverConfig.TokenURL,
	}

	// configure service
	cfg := frontend.ServiceConfig{
		Config:          serverConfig,
		OAuth2Endpoints: oauth2Endpoint,
		CertKeyPair:     keypair,
		Httponly:        httpOnly,
	}
	svc := frontend.Service{}
	err = svc.Run(cfg)
	if err != nil {
		fmt.Printf("ERR: Could not start service, caused by: %s\n", err)
		os.Exit(1)
	}
}
