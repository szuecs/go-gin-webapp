//handles the configuration of the applications. Yaml files are mapped with the struct

package client

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

// ClientConfig is the configuration from the client. Usually loaded from config files.
type ClientConfig struct {
	URL           string   //URL to our service endpoint set by the user
	RealURL       *url.URL //RealURL to our service endpoint parsed from URL
	Debug         bool     //true if Debug is enabled
	Oauth2Enabled bool     //true if oauth2 is enabled
	OauthURL      string   //the oauth2 endpoint to be used
	TokenURL      string   //the oauth2 token info endpoint
	Username      string   //user to authenticate with, to get a token
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
	viper := viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/go-gin-webapp-cli")
	viper.AddConfigPath(fmt.Sprintf("%s/.config/go-gin-webapp-cli", os.ExpandEnv("$HOME")))
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Can not read config, caused by: %s", err)
		return nil, err
	}
	var cfg ClientConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf("Can not marshal config, caused by: %s", err)
		return nil, err
	}

	return &cfg, nil
}
