package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/szuecs/go-gin-webapp/client"
	konfig "github.com/szuecs/go-gin-webapp/conf/client"
)

//Buildstamp and Githash are used to set information at build time regarding
//the version of the build.
//Buildstamp is used for storing the timestamp of the build
var Buildstamp = "Not set"

//Githash is used for storing the commit hash of the build
var Githash = "Not set"

// Version is used to store the tagged version of the build
var Version = "Not set"

func main() {
	var (
		debug       = kingpin.Flag("debug", "enable debug mode").Default("false").Bool()
		username    = kingpin.Flag("username", "Set username to authenticate with.").Default("").String()
		oauth2Token = kingpin.Flag("oauth2-token", "Set OAuth2 Access Token.").Default("").String()
		baseURL     = kingpin.Flag("url", "Set Base URL.").Default("http://127.0.0.1:8080").String()
		_           = kingpin.Command("root", "Handle root ressource.")
		_           = kingpin.Command("health", "Handle health ressource.")
		_           = kingpin.Command("login", "Handle login and save access-token to .gin-oauth-token.")
		_           = kingpin.Command("version", "show version")
	)

	switch kingpin.Parse() {
	case "login":
		cli := createClient(*baseURL, *oauth2Token, *username, *debug)
		cli.RenewAccessToken()
	case "root":
		cli := createClient(*baseURL, *oauth2Token, *username, *debug)
		cli.Config.RealURL.Path = fmt.Sprintf("%s%s", cli.Config.RealURL.Path, "/")
		cli.GetAccessToken()
		cli.Get(cli.Config.RealURL)
	case "health":
		cli := createClient(*baseURL, *oauth2Token, *username, *debug)
		cli.Config.RealURL.Path = fmt.Sprintf("%s%s", cli.Config.RealURL.Path, "/health")
		cli.GetAccessToken()
		cli.Get(cli.Config.RealURL)
	case "version":
		fmt.Printf(`%s Version: %s
================================
    Buildtime: %s
    GitHash: %s
`, path.Base(os.Args[0]), Version, Buildstamp, Githash)
	}
}

func createClient(url, token, username string, debug bool) client.Client {
	//loading cfg from file. it is overridden by the command line parameters
	cfg := konfig.New()

	if debug {
		cfg.Debug = debug
		fmt.Println("Enabled debug mode")
	}

	// URL, cli parameter overwrites config
	if url != "" {
		cfg.URL = url
	} else if cfg.URL == "" {
		// fallback
		cfg.URL = "http://127.0.0.1:8080"
	}
	cfg.RealURL = parseURL(cfg.URL)

	// client
	cli := client.Client{
		Config:      cfg,
		AccessToken: token,
	}

	// username
	cli.GetUsername(username)

	// AccessToken
	if cli.AccessToken == "" {
		cli.GetAccessToken()
	}

	return cli
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Printf("Failed to parse url %s, caused by: %s", s, err)
		os.Exit(2)
	}
	return u
}
