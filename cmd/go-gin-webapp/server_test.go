package main

import (
	"fmt"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/szuecs/go-gin-webapp/conf"
)

var basePath string

func setBasePath() {
	_, filename, _, _ := runtime.Caller(0)
	basePath = path.Dir(path.Dir(path.Dir(filename)))
}

func getValidHTTPConfig() *conf.Config {
	return &conf.Config{
		Oauth2Enabled:    false,
		DebugEnabled:     false,
		Port:             8080,
		MonitorPort:      9000,
		TLSCertfilePath:  "/does/not/exist",
		TLSKeyfilePath:   "/does/not/exist",
		LogFlushInterval: time.Second * 1,
		AuthURL:          "",
		TokenURL:         "",
		AuthorizedTeams:  []conf.AccessTuple{},
		AuthorizedUsers:  []conf.AccessTuple{},
	}
}

func getInvalidHTTPConfig() *conf.Config {
	return &conf.Config{
		Oauth2Enabled:    true,
		DebugEnabled:     false,
		Port:             8080,
		MonitorPort:      9000,
		TLSCertfilePath:  "/does/not/exist",
		TLSKeyfilePath:   "/does/not/exist",
		LogFlushInterval: time.Second * 1,
		AuthURL:          "",
		TokenURL:         "",
		AuthorizedTeams:  []conf.AccessTuple{},
		AuthorizedUsers:  []conf.AccessTuple{},
	}
}

func getValidHTTPSConfig() *conf.Config {
	if basePath == "" {
		setBasePath()
	}
	return &conf.Config{
		Oauth2Enabled:    false,
		DebugEnabled:     false,
		Port:             8080,
		MonitorPort:      9000,
		TLSCertfilePath:  fmt.Sprintf("%s/test/test.crt", basePath),
		TLSKeyfilePath:   fmt.Sprintf("%s/test/test.key", basePath),
		LogFlushInterval: time.Second * 1,
		AuthURL:          "",
		TokenURL:         "",
		AuthorizedTeams:  []conf.AccessTuple{},
		AuthorizedUsers:  []conf.AccessTuple{},
	}
}

func getInvalidHTTPSConfigCert() *conf.Config {
	if basePath == "" {
		setBasePath()
	}
	return &conf.Config{
		Oauth2Enabled:    false,
		DebugEnabled:     false,
		Port:             8080,
		MonitorPort:      9000,
		TLSCertfilePath:  "/does/not/exist",
		TLSKeyfilePath:   fmt.Sprintf("%s/test/test.key", basePath),
		LogFlushInterval: time.Second * 1,
		AuthURL:          "",
		TokenURL:         "",
		AuthorizedTeams:  []conf.AccessTuple{},
		AuthorizedUsers:  []conf.AccessTuple{},
	}
}

func getInvalidHTTPSConfigKey() *conf.Config {
	if basePath == "" {
		setBasePath()
	}
	return &conf.Config{
		Oauth2Enabled:    false,
		DebugEnabled:     false,
		Port:             8080,
		MonitorPort:      9000,
		TLSCertfilePath:  fmt.Sprintf("%s/test/test.crt", basePath),
		TLSKeyfilePath:   "/does/not/exist",
		LogFlushInterval: time.Second * 1,
		AuthURL:          "",
		TokenURL:         "",
		AuthorizedTeams:  []conf.AccessTuple{},
		AuthorizedUsers:  []conf.AccessTuple{},
	}
}

func TestIsHTTPOnly(t *testing.T) {
	var cfg *conf.Config
	cfg = getValidHTTPConfig()
	if !IsHTTPOnly(cfg) {
		t.Fatal("IsHTTPOnly return false on valid http config, but should return true")
	}
	cfg = getValidHTTPSConfig()
	if IsHTTPOnly(cfg) {
		t.Fatal("IsHTTPOnly return true on valid https config, but should return false")
	}
	cfg = getInvalidHTTPConfig()
	if !IsHTTPOnly(cfg) {
		t.Fatal("IsHTTPOnly return false on invalid http config, but should return true")
	}
	cfg = getInvalidHTTPSConfigCert()
	if !IsHTTPOnly(cfg) {
		t.Fatal("IsHTTPOnly return false on invalid https config, but should return true")
	}
	cfg = getInvalidHTTPSConfigKey()
	if !IsHTTPOnly(cfg) {
		t.Fatal("IsHTTPOnly return false on invalid https config, but should return true")
	}
}

func TestGetServiceConfig(t *testing.T) {
	var cfg *conf.Config
	cfg = getValidHTTPConfig()
	svcConfig, err := GetServiceConfig(cfg, true)
	if err != nil {
		t.Fatal("GetServiceConfig should not return an error")
	}
	if svcConfig == nil {
		t.Fatal("GetServiceConfig should not return nil for config")
	}

	cfg = getValidHTTPSConfig()
	svcConfig, err = GetServiceConfig(cfg, false)
	if err != nil {
		t.Fatal("GetServiceConfig should not return an error")
	}
	if svcConfig == nil {
		t.Fatal("GetServiceConfig should not return nil for config")
	}

	cfg = getInvalidHTTPSConfigCert()
	_, err = GetServiceConfig(cfg, false)
	if err == nil {
		t.Fatal("GetServiceConfig should return an error")
	}
	cfg = getInvalidHTTPSConfigKey()
	_, err = GetServiceConfig(cfg, false)
	if err == nil {
		t.Fatal("GetServiceConfig should return an error")
	}
}
