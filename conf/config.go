// Package conf handles the configuration of the applications. Yaml
// files are mapped with the struct
package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v1"

	"github.com/golang/glog"
)

// Config is the configuration struct. The config file config.yaml
// will unmarshaled to this struct.
type Config struct {
	DebugEnabled     bool          `yaml:"debugEnabled,omitempty"`
	Oauth2Enabled    bool          `yaml:"oauth2Enabled,omitempty"`
	ProfilingEnabled bool          `yaml:"profilingEnabled,omitempty"`
	Port             int           `yaml:"port,omitempty"`
	MonitorPort      int           `yaml:"monitorPort,omitempty"`
	LogFlushInterval time.Duration `yaml:"logFlushInterval,omitempty"`
	TLSCertfilePath  string        `yaml:"tlsCertfilePath,omitempty"`
	TLSKeyfilePath   string        `yaml:"tlsKeyfilePath,omitempty"`
	AuthURL          string        `yaml:"authURL,omitempty"`
	TokenURL         string        `yaml:"tokenURL,omitempty"`
	AuthorizedTeams  []AccessTuple `yaml:"AuthorizedTeams,omitempty"`
	AuthorizedUsers  []AccessTuple `yaml:"AuthorizedUsers,omitempty"`
}

// AccessTuple represents an entry for Authorization
type AccessTuple struct {
	Realm string `yaml:"Realm,omitempty"`
	UID   string `yaml:"UID,omitempty"`
	Cn    string `yaml:"Cn,omitempty"`
}

// shared state for configuration
var conf *Config

// New returns the loaded configuration or panic
func New() (*Config, error) {
	var err error
	if conf == nil {
		conf, err = configInit("config.yaml")
	}
	return conf, err
}

// PROJECTNAME TODO: should be replaced in your application
const PROJECTNAME string = "go-gin-webapp"

// FIXME: not windows compatible
func configInit(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/%s/config.yaml", os.ExpandEnv("$HOME"), PROJECTNAME))
	if err != nil {
		glog.Fatalf("Can not read config, caused by: %s", err)
	}
	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		glog.Fatalf("configuration could not be unmarshaled, caused by: %s", err)
	}
	return &config, err
}
