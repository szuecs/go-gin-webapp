package main

import (
	"fmt"
	"net/url"
	"os"
	"os/user"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/szuecs/go-gin-webapp/client"
	konfig "github.com/szuecs/go-gin-webapp/conf/client"
	"github.com/vrischmann/envconfig"
)

//Buildstamp and Githash are used to set information at build time regarding
//the version of the build.
//Buildstamp is used for storing the timestamp of the build
var Buildstamp = "Not set"

//Githash is used for storing the commit hash of the build
var Githash = "Not set"

// Version is used to store the tagged version of the build
var Version = "Not set"

// flag variables
var version bool

// DEBUG enables debug mode in the cli. Used for verbose printing
var DEBUG bool

var conf struct {
	AccessUser     string `envconfig:"optional"`
	AccessPassword string `envconfig:"optional"`
	AccessToken    string `envconfig:"optional"`
	OAuth2Endpoint struct {
		AuthURL      string `envconfig:"optional"`
		TokenInfoURL string `envconfig:"optional"`
	}
}

func main() {

	usage := fmt.Sprintf(`Usage:
  go-gin-webapp-cli -h | --help
  go-gin-webapp-cli --version
  go-gin-webapp-cli root [options]
  go-gin-webapp-cli health [options]
  go-gin-webapp-cli login [<username>] [options]


Options:
  --oauth2  OAuth2 enable
  --oauth2-token=<access_token>  OAuth2 AccessToken (no user, password required)
  --oauth2-authurl=<oauth2_authurl>  OAuth2 endpoint that issue AccessTokens
  --url=<url>  URL of the endpoint you want to request to.
  --debug  Debug
`)

	arguments, err := docopt.Parse(usage, nil, true, fmt.Sprintf("%s Version %s - Build Time: %s - Git Commit Hash: %s", os.Args[0], Version, Buildstamp, Githash), false)
	if err != nil {
		panic("Could not parse CLI")
	}

	DEBUG = arguments["--debug"].(bool)
	// Auth information from ENV and parameter
	if err := envconfig.Init(&conf); err != nil {
		fmt.Printf("ERR: envconfig failed, caused by: %s\n", err)
		os.Exit(2)
	}

	cli := createClient(arguments)

	username := ""
	// maybe user is not configured
	if cli.Config.Username == "" {
		// guess user from system
		user, err := user.Current()
		if err == nil {
			username = user.Username
		}
	} else {
		username = cli.Config.Username
	}
	// user passed by cli parameter will take precedence
	username = GetStringFromArgs(arguments, "<username>", username)
	if username == "" {
		fmt.Println("Can not find user to authenticate with.")
		os.Exit(2)
	}

	if arguments["root"].(bool) {
		cli.Config.RealUrl.Path = fmt.Sprintf("%s%s", cli.Config.RealUrl.Path, "/")
		cli.GetAccessToken(username)
		cli.Get(cli.Config.RealUrl)
	} else if arguments["health"].(bool) {
		cli.Config.RealUrl.Path = fmt.Sprintf("%s%s", cli.Config.RealUrl.Path, "/health")
		cli.GetAccessToken(username)
		cli.Get(cli.Config.RealUrl)
	} else if arguments["login"].(bool) {
		cli.RenewAccessToken(strings.TrimSpace(username))
	}
}

func createClient(arguments map[string]interface{}) client.Client {
	//loading cfg from file. it is overridden by the command line parameters
	cfg := konfig.New()

	_baseUrl := arguments["--url"]
	// cli parameter overwrites config
	if _baseUrl != nil {
		cfg.Url = _baseUrl.(string)
	} else if cfg.Url == "" {
		// fallback
		cfg.Url = "http://127.0.0.1:8080"
	}
	cfg.RealUrl = parseUrl(cfg.Url)

	accessToken := GetStringFromArgs(arguments, "--oauth2-token", "")

	return client.Client{
		Config:      cfg,
		AccessToken: accessToken,
	}
}

func parseUrl(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Printf("Failed to parse url %s, caused by: %s", s, err)
		os.Exit(2)
	}
	return u
}
