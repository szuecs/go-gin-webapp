package frontend

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"gopkg.in/mcuadros/go-monitor.v1/aspects"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/szuecs/go-gin-webapp/conf"
	"github.com/zalando-techmonkeys/gin-glog"
	"github.com/zalando-techmonkeys/gin-gomonitor"
	"github.com/zalando-techmonkeys/gin-oauth2"
	"github.com/zalando-techmonkeys/gin-oauth2/zalando"
	"golang.org/x/oauth2"
)

// ServiceConfig contains everything configurable for the service
// endpoint.
type ServiceConfig struct {
	Config          *conf.Config
	OAuth2Endpoints oauth2.Endpoint
	CertKeyPair     tls.Certificate
	Httponly        bool
}

var cfg conf.Config

// Service is the main struct
type Service struct {
	Healthy bool
}

func (svc *Service) checkDependencies() bool {
	// TODO: you may want to check if you can connect to your dependencies here
	return true
}

// Run is the main function of the server. It bootstraps the service
// and creates the route endpoints.
func (svc *Service) Run(config *ServiceConfig) error {
	cfg = *config.Config

	// init gin
	if !cfg.DebugEnabled {
		gin.SetMode(gin.ReleaseMode)
	}

	// Middleware
	router := gin.New()
	router.Use(ginglog.Logger(cfg.LogFlushInterval))
	router.Use(gomonitor.Metrics(cfg.MonitorPort, []aspects.Aspect{}))
	router.Use(gin.Recovery())

	// OAuth2 secured if conf.Oauth2Enabled is set
	var private *gin.RouterGroup
	if cfg.Oauth2Enabled {
		zalando.AccessTuples = []zalando.AccessTuple{{"employees", "sszuecs", "Sandor Szücs"}, {"employees", "njuettner", "Nick Jüttner"}, {"employees", "rdifazio", "Raffaele Di Fazio"}}
		private = router.Group("")
		private.Use(ginoauth2.Auth(zalando.UidCheck, config.OAuth2Endpoints))
	}

	//
	//  Handlers
	//
	router.GET("/health", svc.HealthHandler)
	if cfg.Oauth2Enabled {
		// authenticated and authorized routes
		private.GET("/", svc.RootHandler)
	} else {
		// public routes
		router.GET("/", svc.RootHandler)
	}

	// TLS config
	tlsConfig := tls.Config{}
	if !config.Httponly {
		tlsConfig.Certificates = []tls.Certificate{config.CertKeyPair}
		tlsConfig.NextProtos = []string{"http/1.1"}
		tlsConfig.Rand = rand.Reader // Strictly not necessary, should be default
	}

	// run frontend server
	serve := &http.Server{
		Addr:      fmt.Sprintf(":%d", cfg.Port),
		Handler:   router,
		TLSConfig: &tlsConfig,
	}

	if cfg.ProfilingEnabled {
		ginpprof.Wrapper(router)
	}

	if svc.checkDependencies() {
		svc.Healthy = true
	}

	// start server
	if config.Httponly {
		err := serve.ListenAndServe()
		if err != nil {
			glog.Fatalf("Can not Serve HTTP, caused by: %s\n", err)
		}
	} else {
		conn, err := net.Listen("tcp", serve.Addr)
		if err != nil {
			panic(err)
		}
		tlsListener := tls.NewListener(conn, &tlsConfig)
		err = serve.Serve(tlsListener)
		if err != nil {
			glog.Fatalf("Can not Serve TLS, caused by: %s\n", err)
		}
	}
	return nil
}
