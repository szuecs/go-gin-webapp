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
	DebugEnabled     bool          `yaml:"debug_enabled,omitempty"`
	Oauth2Enabled    bool          `yaml:"oauth2_enabled,omitempty"`
	ProfilingEnabled bool          `yaml:"profiling_enabled,omitempty"`
	Port             int           `yaml:"port,omitempty"`
	MonitorPort      int           `yaml:"monitor_port,omitempty"`
	LogFlushInterval time.Duration `yaml:"log_flush_interval,omitempty"`
	TLSCertfilePath  string        `yaml:"tls_certfile_path,omitempty"`
	TLSKeyfilePath   string        `yaml:"tls_keyfile_path,omitempty"`
	AuthURL          string        `yaml:"auth_url,omitempty"`
	TokenURL         string        `yaml:"token_url,omitempty"`
	AuthorizedTeams  []AccessTuple `yaml:"authorized_teams,omitempty"`
	AuthorizedUsers  []AccessTuple `yaml:"authorized_users,omitempty"`
}

// AccessTuple represents an entry for Authorization
type AccessTuple struct {
	Realm string `yaml:"realm,omitempty"`
	UID   string `yaml:"uid,omitempty"`
	Cn    string `yaml:"cn,omitempty"`
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
func readFile(filepath string) ([]byte, bool) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return b, false
	}
	return b, true
}

// FIXME: not windows compatible
func configInit(filename string) (*Config, error) {
	globalConfig := fmt.Sprintf("/etc/%s/config.yaml", PROJECTNAME)
	homeConfig := fmt.Sprintf("%s/.config/%s/config.yaml", os.ExpandEnv("$HOME"), PROJECTNAME)
	b, ok := readFile(homeConfig)
	if !ok {
		b, ok = readFile(globalConfig)
	}
	if !ok {
		return nil, fmt.Errorf("No file readable in %v nor in %v", globalConfig, homeConfig)
	}
	var config Config
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		glog.Fatalf("configuration could not be unmarshaled, caused by: %s", err)
	}
	return &config, err
}
