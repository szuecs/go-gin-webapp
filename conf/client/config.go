//handles the configuration of the applications. Yaml files are mapped with the struct

package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"gopkg.in/yaml.v1"

	"github.com/golang/glog"
)

// ClientConfig is the configuration from the client. Usually loaded from config files.
type ClientConfig struct {
	URL           string   `yaml:"URL,omitempty"` //URL to our service endpoint set by the user
	RealURL       *url.URL //RealURL to our service endpoint parsed from URL
	Debug         bool     `yaml:"debugEnabled,omitempty"`  //true if Debug is enabled
	Oauth2Enabled bool     `yaml:"oauth2Enabled,omitempty"` //true if oauth2 is enabled
	OauthURL      string   `yaml:"URL,omitempty"`           //the oauth2 endpoint to be used
	TokenURL      string   `yaml:"URL,omitempty"`           //the oauth2 token info endpoint
	Username      string   `yaml:"URL,omitempty"`           //user to authenticate with, to get a token
}

//shared state for configuration
var clientConf *ClientConfig

//New gets the ClientConfiguration
func New() *ClientConfig {
	if clientConf == nil {
		var err error
		clientConf, err = clientConfigInit("config.yaml")
		if err != nil {
			fmt.Printf("could not load configuration. Reason: %s\n", err)
			os.Exit(2)
		}
	}

	return clientConf
}

//FIXME: ConfigPath and expansion of ENV vars like this is not windows compatible
func clientConfigInit(filename string) (*ClientConfig, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/go-gin-webapp-cli/config.yaml", os.ExpandEnv("$HOME")))
	if err != nil {
		glog.Fatalf("Can not read config, caused by: %s", err)
	}
	var cfg ClientConfig
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		glog.Fatalf("configuration could not be unmarshaled, caused by: %s", err)
	}
	return &cfg, err
}
