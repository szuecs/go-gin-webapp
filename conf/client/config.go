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
	URL           string   `yaml:"url,omitempty"` //URL to our service endpoint set by the user
	RealURL       *url.URL //RealURL to our service endpoint parsed from URL
	Debug         bool     `yaml:"debug_enabled,omitempty"`  //true if Debug is enabled
	Oauth2Enabled bool     `yaml:"oauth2_enabled,omitempty"` //true if oauth2 is enabled
	OauthURL      string   `yaml:"oauth2_url,omitempty"`     //the oauth2 endpoint to be used
	TokenURL      string   `yaml:"token_url,omitempty"`      //the oauth2 token info endpoint
	Username      string   `yaml:"username,omitempty"`       //user to authenticate with, to get a token
}

//shared state for configuration
var clientConf *ClientConfig

//New gets the ClientConfiguration
func New() *ClientConfig {
	if clientConf == nil {
		var err error
		clientConf, err = clientConfigInit("config.yaml")
		if err != nil {
			glog.Exitf("could not load configuration. Reason: %s\n", err)
		}
	}

	return clientConf
}

//FIXME: ConfigPath and expansion of ENV vars like this is not windows compatible
func clientConfigInit(filename string) (*ClientConfig, error) {
	fpath := fmt.Sprintf("%s/.config/go-gin-webapp-cli/config.yaml", os.ExpandEnv("$HOME"))
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		glog.Exitf("can not read config from %s, caused by: %s", fpath, err)
	}
	var cfg ClientConfig
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		glog.Exitf("configuration could not be unmarshaled, caused by: %s", err)
	}
	return &cfg, err
}
